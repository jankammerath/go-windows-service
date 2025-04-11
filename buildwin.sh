#!/bin/sh

# Set default architecture to arm64
ARCH="arm64"

# Check if an architecture is provided as an argument
if [ -n "$1" ]; then
  case "$1" in
    amd64)
      ARCH="amd64"
      ;;
    arm64)
      ARCH="arm64"
      ;;
    *)
      echo "Invalid architecture: $1. Using default: arm64"
      ;;
  esac
fi

rm -r bin
mkdir -p bin

# Build the Windows executable
echo "Building Windows executable..."
GOOS=windows GOARCH="$ARCH" CGO_ENABLED=0 go build -o bin/cpuservice.exe -ldflags '-w -s'

# Check if build was successful
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

# Create a ZIP file with the executable and service scripts
echo "Creating distribution ZIP file..."
ZIP_NAME="cpuservice-$ARCH.zip"

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