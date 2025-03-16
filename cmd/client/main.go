package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const PORT = ":8080"

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
		sendText(serverIP, strings.Join(os.Args[3:], " "))

	case "file":
		if len(os.Args) < 4 {
			fmt.Println("Error: File path is required")
			return
		}
		sendFile(serverIP, os.Args[3])

	default:
		fmt.Println("Unknown command. Use 'text' or 'file'")
	}
}

func sendText(serverIP, message string) {
	conn, err := net.Dial("tcp", serverIP+PORT)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer conn.Close()

	// Send the message
	writer := bufio.NewWriter(conn)
	_, err = writer.WriteString(message + "\n")
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return
	}

	// Ensure the message is sent
	if err := writer.Flush(); err != nil {
		fmt.Printf("Error flushing message: %v\n", err)
		return
	}

	fmt.Println("Message sent successfully")
}

func sendFile(serverIP, filePath string) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Connect to server
	conn, err := net.Dial("tcp", serverIP+PORT)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer conn.Close()

	// Send the filename
	writer := bufio.NewWriter(conn)
	filename := filepath.Base(filePath)
	_, err = writer.WriteString("FILE:" + filename + "\n")
	if err != nil {
		fmt.Printf("Error sending filename: %v\n", err)
		return
	}
	writer.Flush()

	// Send the file content
	_, err = io.Copy(conn, file)
	if err != nil {
		fmt.Printf("Error sending file: %v\n", err)
		return
	}

	fmt.Printf("File %s sent successfully\n", filename)
}
