package main

import (
	"fmt"
	"os"
	"path/filepath"

	"local-share/pkg/receiver"
	"local-share/pkg/sender"
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
		receiver.Start()
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
			sender.SendText(serverIP, message)
		case "file":
			// Check arguments
			if len(os.Args) < 5 {
				fmt.Println("Usage: local-share send file <server-ip> <filepath>")
				os.Exit(1)
			}

			serverIP := os.Args[3]
			filePath := os.Args[4]

			// Run the client file sending functionality
			sender.SendFile(serverIP, filePath)
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
