# Go Windows Service

A basic Windows Service that listens on port `8899` providing HTTP serving CPU stats as `json`.

## Installation

### Easy Installation (Recommended)

1. Extract the ZIP file to a location of your choice
2. Right-click on `install-service.bat` and select **Run as administrator**
3. The script will automatically:
   - Install the service as "Go CPU Monitor Service"
   - Configure it to start automatically
   - Start the service immediately

### Uninstalling the Service

1. Right-click on `uninstall-service.bat` and select **Run as administrator**
2. The script will automatically:
   - Stop the running service
   - Remove it from the Windows services list

### Manual Installation (Alternative)

If you prefer to install the service manually using the Service Control Manager:

1. Copy the `cpuservice.exe` file to the desired location
2. Open Command Prompt as Administrator
3. Run: `sc.exe create GoCPUMonitor binPath= "C:\path\to\cpuservice.exe" start= auto DisplayName= "Go CPU Monitor Service"`
4. Run: `sc.exe start GoCPUMonitor`

Note: The space after `=` is required by the sc.exe command syntax

## Accessing the Service

Once the service is running, you can access CPU statistics via:

```
http://localhost:8899/cpu
```

The response will be in JSON format:

```json
{
  "utilization": 23.45
}
```

## Troubleshooting

If the service fails to start:

1. Check the Windows Event Viewer for errors
2. Verify that port 8899 is not in use by another application
3. Ensure you're running the installation scripts as Administrator

