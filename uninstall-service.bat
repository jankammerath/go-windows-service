@echo off
echo Uninstalling Go CPU Monitor Service...

set SERVICE_NAME=GoCPUMonitor

REM Check if running as administrator
net session >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: This script requires administrator privileges.
    echo Please right-click this batch file and select "Run as administrator"
    pause
    exit /b 1
)

REM Check if service exists
sc query %SERVICE_NAME% >nul 2>&1
if %errorlevel% neq 0 (
    echo Service "%SERVICE_NAME%" does not exist or is already uninstalled.
    pause
    exit /b 0
)

REM Stop the service first
echo Stopping service...
sc stop %SERVICE_NAME%
timeout /t 2 >nul

REM Delete the service
echo Deleting service...
sc delete %SERVICE_NAME%
if %errorlevel% neq 0 (
    echo Failed to uninstall service "%SERVICE_NAME%".
    pause
    exit /b 1
)

echo.
echo Service "%SERVICE_NAME%" uninstalled successfully.
echo.
pause