# Local Share

A simple CLI application written in Go that allows sharing text and files between computers in a Local Area Network (LAN).

## Features

- Share text messages between computers
- Share files between computers
- Works on any computer in the same LAN
- Simple command-line interface

## Prerequisites

- Go 1.21 or later installed on your system

## Project Structure

```
local-share/
├── bin/            # Compiled binaries
├── cmd/
│   ├── server/     # Server application
│   └── client/     # Client application
├── uploads/        # Directory for received files
└── go.mod
```

## Building

To build the application, run these commands from the project root:

For Windows:
```bash
# Create bin directory if it doesn't exist
mkdir bin

# Build server
go build -o bin/server.exe ./cmd/server

# Build client
go build -o bin/client.exe ./cmd/client
```

For Unix-like systems (Linux/macOS):
```bash
# Create bin directory if it doesn't exist
mkdir -p bin

# Build server
go build -o bin/server ./cmd/server

# Build client
go build -o bin/client ./cmd/client
```

This will create two executables in the `bin` directory:
- `bin/server` (on Unix-like systems) or `bin/server.exe` (on Windows)
- `bin/client` (on Unix-like systems) or `bin/client.exe` (on Windows)

## Usage

### Starting the Server

1. Run the server on the receiving computer:
```bash
# Windows
.\bin\server.exe

# Unix-like systems
./bin/server
```

The server will display its IP address and start listening on port 8080.

### Sending Text

To send text to the server, use:
```bash
# Windows
.\bin\client.exe text <server-ip> <message>

# Unix-like systems
./bin/client text <server-ip> <message>
```

Example:
```bash
# Windows
.\bin\client.exe text 192.168.1.100 Hello, this is a test message!

# Unix-like systems
./bin/client text 192.168.1.100 Hello, this is a test message!
```

### Sending Files

To send a file to the server, use:
```bash
# Windows
.\bin\client.exe file <server-ip> <file-path>

# Unix-like systems
./bin/client file <server-ip> <file-path>
```

Example:
```bash
# Windows
.\bin\client.exe file 192.168.1.100 C:\path\to\your\file.txt

# Unix-like systems
./bin/client file 192.168.1.100 /path/to/your/file.txt
```

## Notes

- The server creates an `uploads` directory to store received files
- Make sure both computers are on the same network
- The server's IP address is displayed when you start it
- Port 8080 must be available on the server 