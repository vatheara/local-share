package main

import (
	"fmt"
	"os"
	"strings"

	"local-share/client"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  To send text: local-share text <server-ip> <message>")
		fmt.Println("  To send file: local-share file <server-ip> <file-path>")
		return
	}

	command := os.Args[1]
	serverIP := os.Args[2]

	switch command {
	case "text":
		if len(os.Args) < 4 {
			fmt.Println("Error: Message is required")
			return
		}
		client.SendText(serverIP, strings.Join(os.Args[3:], " "))

	case "file":
		if len(os.Args) < 4 {
			fmt.Println("Error: File path is required")
			return
		}
		client.SendFile(serverIP, os.Args[3])

	default:
		fmt.Println("Unknown command. Use 'text' or 'file'")
	}
}
