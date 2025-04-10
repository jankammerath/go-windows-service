#!/bin/sh
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/cpuservice.exe -ldflags '-w -s'