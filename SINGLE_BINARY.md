# Converting to a Single Binary Approach

## Overview

Instead of having separate binaries for the server and client, we'll create a single binary that supports different commands:

```bash
# Start the server
local-share receiver

# Send text
local-share send text <server-ip> "Your message"

# Send file
local-share send file <server-ip> /path/to/file
```

## Code Structure Changes

### 1. Create a new main.go file in the cmd directory

```go
// cmd/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "receiver":
		// Run the server functionality
		runReceiver()
	case "send":
		if len(os.Args) < 3 {
			printUsage()
			os.Exit(1)
		}
		
		subCommand := os.Args[2]
		
		switch subCommand {
		case "text":
			// Check arguments
			if len(os.Args) < 5 {
				fmt.Println("Usage: local-share send text <server-ip> <message>")
				os.Exit(1)
			}
			
			serverIP := os.Args[3]
			message := os.Args[4]
			
			// Run the client text sending functionality
			runSendText(serverIP, message)
		case "file":
			// Check arguments
			if len(os.Args) < 5 {
				fmt.Println("Usage: local-share send file <server-ip> <filepath>")
				os.Exit(1)
			}
			
			serverIP := os.Args[3]
			filePath := os.Args[4]
			
			// Run the client file sending functionality
			runSendFile(serverIP, filePath)
		default:
			fmt.Printf("Unknown send subcommand: %s\n", subCommand)
			printUsage()
			os.Exit(1)
		}
	case "--help", "-h", "help":
		printUsage()
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf("Usage: %s COMMAND [ARGS...]\n\n", filepath.Base(os.Args[0]))
	fmt.Println("Commands:")
	fmt.Println("  receiver                  Start the receiver server")
	fmt.Println("  send text <ip> <message>  Send a text message to a server")
	fmt.Println("  send file <ip> <filepath> Send a file to a server")
	fmt.Println("  help                      Show this help message")
}

// These functions will import and call the actual functionality from the existing code
func runReceiver() {
	// Import and call the server code
	// You'll need to refactor the server code to be importable as a package
}

func runSendText(serverIP, message string) {
	// Import and call the client text sending code
	// You'll need to refactor the client code to be importable as a package
}

func runSendFile(serverIP, filePath string) {
	// Import and call the client file sending code
	// You'll need to refactor the client code to be importable as a package
}
```

### 2. Refactor Server Code

Move the server code from `cmd/server/main.go` to a package:

```go
// cmd/receiver/receiver.go
package receiver

import (
    // existing imports
)

// Export the main functionality as functions
func Start() {
    // The code that was in main() in the server
}

// Export other functions as needed
```

### 3. Refactor Client Code

Move the client code from `cmd/client/main.go` to a package:

```go
// cmd/sender/sender.go
package sender

import (
    // existing imports
)

// Export the text sending functionality
func SendText(serverIP, message string) {
    // The code that was handling text sending in the client
}

// Export the file sending functionality
func SendFile(serverIP, filePath string) {
    // The code that was handling file sending in the client
}
```

### 4. Update go.mod

The module structure remains the same, but you'll be organizing the code differently.

### 5. Update the build commands

Replace the separate build commands with:

```bash
go build -o bin/local-share ./cmd
```

## Testing

Make sure to test all the commands thoroughly after refactoring:

```bash
# Start the receiver
./bin/local-share receiver

# Send text
./bin/local-share send text 192.168.1.100 "Test message"

# Send file
./bin/local-share send file 192.168.1.100 ./test-file.txt
```

## Updating the Homebrew Formula

After implementing the single binary approach, the Homebrew formula will be simpler as shown in the updated `local-share.rb` file. 