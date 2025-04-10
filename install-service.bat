@echo off
echo Installing Go CPU Monitor Service...

REM Get the directory of this batch file
set SCRIPT_DIR=%~dp0
set SERVICE_NAME=GoCPUMonitor
set EXECUTABLE_PATH=%SCRIPT_DIR%bin\cpuservice.exe

echo Service executable path: %EXECUTABLE_PATH%

REM Check if running as administrator
net session >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: This script requires administrator privileges.
    echo Please right-click this batch file and select "Run as administrator"
    pause
    exit /b 1
)

REM Check if service already exists
sc query %SERVICE_NAME% >nul 2>&1
if %errorlevel% equ 0 (
    echo Service "%SERVICE_NAME%" already exists.
    echo Stopping and removing existing service...
    sc stop %SERVICE_NAME% >nul 2>&1
    sc delete %SERVICE_NAME% >nul 2>&1
    timeout /t 2 >nul
)

REM Create the service
echo Creating service...
sc create %SERVICE_NAME% binPath= "%EXECUTABLE_PATH%" start= auto DisplayName= "Go CPU Monitor Service"
if %errorlevel% neq 0 (
    echo Failed to create service "%SERVICE_NAME%".
    pause
    exit /b 1
)

REM Start the service
echo Starting service...
sc start %SERVICE_NAME%
if %errorlevel% neq 0 (
    echo Warning: Failed to start service. You may need to start it manually.
    echo You can start it using: sc start %SERVICE_NAME%
) else (
    echo Service started successfully.
)

echo.
echo Service "%SERVICE_NAME%" installed successfully.
echo You can access the CPU data at: http://localhost:8899/cpu
echo.
pause