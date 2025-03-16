# Local Share

A simple CLI application written in Go that allows sharing text and files between computers in a Local Area Network (LAN).

## Features

- Share text messages between computers
- Share files between computers
- Works on any computer in the same LAN
- Simple command-line interface

## Prerequisites

- Go 1.21 or later installed on your system

## Building

To build the application, run:

```bash
go build
```

This will create two executables:
- `local-share` (on Unix-like systems)
- `local-share.exe` (on Windows)

## Usage

### Starting the Server

1. Run the server on the receiving computer:
```bash
./local-share
```

The server will display its IP address and start listening on port 8080.

### Sending Text

To send text to the server, use:
```bash
./local-share text <server-ip> <message>
```

Example:
```bash
./local-share text 192.168.1.100 Hello, this is a test message!
```

### Sending Files

To send a file to the server, use:
```bash
./local-share file <server-ip> <file-path>
```

Example:
```bash
./local-share file 192.168.1.100 /path/to/your/file.txt
```

## Notes

- The server creates an `uploads` directory to store received files
- Make sure both computers are on the same network
- The server's IP address is displayed when you start it
- Port 8080 must be available on the server 