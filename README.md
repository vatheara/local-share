# Local Share

A simple CLI application written in Go that allows sharing text and files between computers in a Local Area Network (LAN).

## Features

- Share text messages between computers
- Share files between computers
- End-to-end encryption for all transfers (AES-256)
  - Encrypted text messages
  - Encrypted file transfers (both filename and content)
- Works on any computer in the same LAN
- Simple command-line interface

## Prerequisites

- Go 1.21 or later installed on your system
- For password input features: `golang.org/x/term` package

## Project Structure

```
local-share/
├── bin/            # Compiled binaries
├── cmd/            # Main application entry point
├── pkg/
│   ├── receiver/   # Server functionality
│   ├── sender/     # Client functionality
│   └── crypto/     # Shared encryption utilities
├── uploads/        # Directory for received files
└── go.mod
```

## Building

To build the application, run these commands from the project root:

For Windows:
```bash
# Create bin directory if it doesn't exist
mkdir bin

# Build the single binary
go build -o bin/local-share.exe ./cmd
```

For Unix-like systems (Linux/macOS):
```bash
# Create bin directory if it doesn't exist
mkdir -p bin

# Build the single binary
go build -o bin/local-share ./cmd
```

This will create a single executable in the `bin` directory:
- `bin/local-share` (on Unix-like systems) or `bin/local-share.exe` (on Windows)

## Usage

### Starting the Server

1. Run the server on the receiving computer:
```bash
# Windows
.\bin\local-share.exe receiver

# Unix-like systems
./bin/local-share receiver
```

The server will prompt for a password to encrypt/decrypt transfers, then display its IP address and start listening on port 8080.

### Sending Text (Encrypted)

To send encrypted text to the server, use:
```bash
# Windows
.\bin\local-share.exe send text <server-ip> <message>

# Unix-like systems
./bin/local-share send text <server-ip> <message>
```

Example:
```bash
# Windows
.\bin\local-share.exe send text 192.168.1.100 "Hello, this is a test message!"

# Unix-like systems
./bin/local-share send text 192.168.1.100 "Hello, this is a test message!"
```

When sending text, you'll be prompted to enter the same password that was used to start the server.

### Sending Files (Encrypted)

To send an encrypted file to the server, use:
```bash
# Windows
.\bin\local-share.exe send file <server-ip> <file-path>

# Unix-like systems
./bin/local-share send file <server-ip> <file-path>
```

Example:
```bash
# Windows
.\bin\local-share.exe send file 192.168.1.100 C:\path\to\your\file.txt

# Unix-like systems
./bin/local-share send file 192.168.1.100 /path/to/your/file.txt
```

The file will be encrypted before transfer, including both the filename and content. The server will decrypt it automatically using the same password.

### Getting Help

To show usage information:
```bash
# Windows
.\bin\local-share.exe help

# Unix-like systems
./bin/local-share help
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
- The server will prompt for a password when starting (with confirmation)
- The client will prompt for the same password when sending messages or files

Important security notes:
- Use the same password on both client and server
- Share the password securely with the receiver (not over the same network)
- Consider changing the password periodically for better security
- The password is automatically converted to a secure encryption key
- All data (text messages, filenames, and file contents) is encrypted
- Even if someone captures the network traffic, they cannot read the data without the password

## Notes

- The server creates an `uploads` directory to store received files
- Make sure both computers are on the same network
- The server's IP address is displayed when you start it
- Port 8080 must be available on the server 