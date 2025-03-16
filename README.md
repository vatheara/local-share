# Local Share

A simple CLI application written in Go that allows sharing text and files between computers in a Local Area Network (LAN).

## Features

- Share text messages between computers
- Share files between computers
- Password-based encryption for text messages (AES-256)
- Works on any computer in the same LAN
- Simple command-line interface

## Prerequisites

- Go 1.21 or later installed on your system
- For password input features: `golang.org/x/term` package

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

The server will prompt for a password to encrypt messages, then display its IP address and start listening on port 8080.

### Sending Text (Encrypted)

To send encrypted text to the server, use:
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

When sending text, you'll be prompted to enter the same password that was used to start the server.

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

## Password Management

The password can be provided in two ways:

1. Environment Variable:
```bash
# Windows PowerShell
$env:LOCALSHARE_KEY="your-password"

# Unix-like systems
export LOCALSHARE_KEY="your-password"
```

2. Interactive Prompt:
- The server will prompt for a password when starting
- The client will prompt for the same password when sending messages

Important security notes:
- Use the same password on both client and server
- Share the password securely with the receiver (not over the same network)
- Consider changing the password periodically for better security
- The password is automatically converted to a secure encryption key

## Notes

- The server creates an `uploads` directory to store received files
- Make sure both computers are on the same network
- The server's IP address is displayed when you start it
- Port 8080 must be available on the server 