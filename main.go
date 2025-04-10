//go:build windows

// --- Installation Instructions (Save as install.bat) ---
// @echo off
// set SERVICE_NAME=GoCPUMonitor
// set EXECUTABLE_PATH=%~dp0%YOUR_GO_EXECUTABLE_NAME.exe
//
// sc.exe create %SERVICE_NAME% binPath= "%EXECUTABLE_PATH%" start= auto DisplayName= "Go CPU Monitor Service"
// if %errorlevel% == 0 (
//     echo Service "%SERVICE_NAME%" created successfully.
//     echo You can start the service using:
//     echo sc.exe start %SERVICE_NAME%
// ) else (
//     echo Failed to create service "%SERVICE_NAME%".
// )
// pause

// --- Uninstallation Instructions (Save as uninstall.bat) ---
// @echo off
// set SERVICE_NAME=GoCPUMonitor
//
// sc.exe stop %SERVICE_NAME%
// sc.exe delete %SERVICE_NAME%
// if %errorlevel% == 0 (
//     echo Service "%SERVICE_NAME%" uninstalled successfully.
// ) else (
//     echo Failed to uninstall service "%SERVICE_NAME%".
// )
// pause

package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Windows API constants and types
type SYSTEM_INFO struct {
	ProcessorArchitecture     uint16
	Reserved                  uint16
	PageSize                  uint32
	MinimumApplicationAddress uintptr
	MaximumApplicationAddress uintptr
	ActiveProcessorMask       uintptr
	NumberOfProcessors        uint32
	ProcessorType             uint32
	AllocationGranularity     uint32
	ProcessorLevel            uint16
	ProcessorRevision         uint16
}

// GetSystemInfo gets information about the system
var (
	modkernel32                    = windows.NewLazySystemDLL("kernel32.dll")
	procGetSystemInfo              = modkernel32.NewProc("GetSystemInfo")
	procGetSystemTimes             = modkernel32.NewProc("GetSystemTimes")
	procGetSystemTimeAsFileTime    = modkernel32.NewProc("GetSystemTimeAsFileTime")
	procGetCurrentProcess          = modkernel32.NewProc("GetCurrentProcess")
	advapi32                       = windows.NewLazySystemDLL("advapi32.dll")
	procStartServiceCtrlDispatcher = advapi32.NewProc("StartServiceCtrlDispatcherW")
	isWindowsService               = true
)

func getSystemInfo() SYSTEM_INFO {
	var si SYSTEM_INFO
	procGetSystemInfo.Call(uintptr(unsafe.Pointer(&si)))
	return si
}

func getSystemTimes(idleTime, kernelTime, userTime *windows.Filetime) error {
	r1, _, e1 := procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(idleTime)),
		uintptr(unsafe.Pointer(kernelTime)),
		uintptr(unsafe.Pointer(userTime)),
	)
	if r1 == 0 {
		return e1
	}
	return nil
}

// CPUUsage struct to hold the CPU utilization
type CPUUsage struct {
	Utilization float64 `json:"utilization"`
}

var (
	lastIdleTicks   int64
	lastSystemTicks int64
	mu              sync.Mutex
)

func getCPUUtilization() float64 {
	mu.Lock()
	defer mu.Unlock()

	// Get number of processors
	sysInfo := getSystemInfo()
	numProcessors := int(sysInfo.NumberOfProcessors)
	if numProcessors == 0 {
		numProcessors = 1 // Fallback if we couldn't get the processor count
	}

	var idle, kernel, user windows.Filetime
	err := getSystemTimes(&idle, &kernel, &user)
	if err != nil {
		log.Printf("Error getting system times: %v", err)
		return 0.0
	}

	// Convert file times to int64 (100-nanosecond intervals)
	idleTime := int64(idle.HighDateTime)<<32 | int64(idle.LowDateTime)
	kernelTime := int64(kernel.HighDateTime)<<32 | int64(kernel.LowDateTime)
	userTime := int64(user.HighDateTime)<<32 | int64(user.LowDateTime)

	// System time is kernel + user time
	// Note: kernel time includes idle time
	systemTime := kernelTime + userTime - idleTime

	idleDiff := idleTime - lastIdleTicks
	systemDiff := systemTime - lastSystemTicks

	lastIdleTicks = idleTime
	lastSystemTicks = systemTime

	if systemDiff <= 0 {
		return 0.0 // Avoid division by zero or negative values
	}

	// Calculate CPU usage percentage across all processors
	cpuUsage := 100.0 * (float64(systemDiff) / (float64(systemDiff) + float64(idleDiff)))

	// Ensure values are within reasonable range
	if cpuUsage < 0 {
		cpuUsage = 0
	} else if cpuUsage > 100.0 {
		cpuUsage = 100.0
	}

	return cpuUsage
}

func cpuHandler(w http.ResponseWriter, r *http.Request) {
	utilization := getCPUUtilization()
	cpuData := CPUUsage{Utilization: utilization}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cpuData)
}

func main() {
	// Check for the "-run" flag to run as a normal application instead of a service
	runAsApp := flag.Bool("run", false, "Run as a normal application instead of a Windows service")
	flag.Parse()

	if *runAsApp {
		isWindowsService = false
		log.Println("Running as a normal application (not a Windows service)")
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	// Initialize last CPU ticks
	getCPUUtilization()
	time.Sleep(time.Millisecond * 100) // Give it a small delay to get initial values
	getCPUUtilization()                // Update again after the delay

	http.HandleFunc("/cpu", cpuHandler)

	port := 8899
	log.Printf("Starting HTTP server on port %d", port)
	server := &http.Server{Addr: ":" + strconv.Itoa(port)}

	if isWindowsService {
		runService(server)
	} else {
		// Run as a regular application for testing
		go func() {
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatalf("HTTP server ListenAndServe error: %v", err)
			}
		}()

		// Handle graceful shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")

		// Add a timeout to the shutdown process
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
		log.Println("Server gracefully stopped.")
	}
}
