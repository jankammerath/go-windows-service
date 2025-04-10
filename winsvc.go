//go:build windows

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

var elog debug.Log

type service struct {
	server *http.Server
}

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (bool, uint32) {
	elog.Info(1, "Go CPU Monitor Service started")
	defer elog.Info(1, "Go CPU Monitor Service stopped")

	changes <- svc.Status{State: svc.StartPending}

	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			elog.Error(1, fmt.Sprintf("HTTP server ListenAndServe error: %v", err))
		}
	}()

	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop}

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Stop:
				elog.Info(1, "Go CPU Monitor Service stopping")
				changes <- svc.Status{State: svc.StopPending}
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := s.server.Shutdown(ctx); err != nil {
					elog.Error(1, fmt.Sprintf("Server shutdown error: %v", err))
				}
				break loop
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// Log current status
				_, err := os.OpenFile("C:\\Temp\\service_status.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err == nil {
					// You could log more detailed status here if needed
				}
			default:
				elog.Error(1, fmt.Sprintf("Unexpected control request #%d", c))
			}
		}
	}

	changes <- svc.Status{State: svc.Stopped}
	return false, 0 // false means no error, 0 is the Windows error code
}

func runService(server *http.Server) {
	var err error
	elog, err = eventlog.Open("GoCPUMonitor")
	if err != nil {
		return
	}
	defer elog.Close()

	elog.Info(1, "Starting Go CPU Monitor Service")
	run := svc.Run
	err = run("GoCPUMonitor", &service{server: server})
	if err != nil {
		elog.Error(1, fmt.Sprintf("Service start failed: %v", err))
		return
	}
	elog.Info(1, "Service stopped successfully")
}
