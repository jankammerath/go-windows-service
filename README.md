# Go Windows Service

A basic Windows Service that listens on port `8899` providing HTTP serving CPU stats as `json`.

## Install the service

### Using sc.exe (Service Control Manager)

1. Copy the `cpuservice.exe` file to the desired location (e.g., `C:\Program Files\CPUMonitor\cpuservice.exe`)

2. Open Command Prompt as Administrator

3. Create the service using sc.exe:

```cmd
sc.exe create GoCPUMonitor binPath= "C:\Program Files\CPUMonitor\cpuservice.exe" start= auto DisplayName= "Go CPU Monitor Service"
```

Note: The space after `=` is required by the sc.exe command syntax

4. Start the service:

```cmd
sc.exe start GoCPUMonitor
```

5. Verify the service is running:

```cmd
sc.exe query GoCPUMonitor
```

### Uninstalling the service

1. Stop the service:

```cmd
sc.exe stop GoCPUMonitor
```

2. Delete the service:

```cmd
sc.exe delete GoCPUMonitor
```

## Accessing the service

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

