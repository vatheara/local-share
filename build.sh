#!/bin/bash

# Script to build local-share for Unix-like systems

# Create bin directory if it doesn't exist
mkdir -p bin

# Build the single binary
echo "Building local-share..."
go build -o bin/local-share ./cmd

if [ $? -eq 0 ]; then
  echo "Build successful! Binary located at bin/local-share"
  echo ""
  echo "Usage examples:"
  echo "  Start server:       ./bin/local-share receiver"
  echo "  Send text message:  ./bin/local-share send text <server-ip> \"message\""
  echo "  Send file:          ./bin/local-share send file <server-ip> /path/to/file"
  echo "  Show help:          ./bin/local-share help"
else
  echo "Build failed!"
fi 