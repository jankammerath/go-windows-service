#!/bin/sh
rm -r bin
mkdir -p bin

# Build the Windows executable
echo "Building Windows executable..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/cpuservice.exe -ldflags '-w -s'

# Check if build was successful
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

# Create a ZIP file with the executable and service scripts
echo "Creating distribution ZIP file..."
ZIP_NAME="cpuservice.zip"

# Check if zip command is available
if ! command -v zip &> /dev/null; then
    echo "The 'zip' command is not installed. Please install it."
    exit 1
fi

# Create the ZIP file
zip -j "bin/$ZIP_NAME" bin/cpuservice.exe install-service.bat uninstall-service.bat

# Check if zip was successful
if [ $? -ne 0 ]; then
    echo "Failed to create ZIP file!"
    exit 1
fi

echo "Build successful! Distribution ZIP created at: bin/$ZIP_NAME"