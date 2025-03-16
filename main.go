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

const (
	PORT        = ":8080"
	BUFFER_SIZE = 1024 * 1024 // 1MB buffer for file transfers
)

func main() {
	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("uploads", 0755); err != nil {
		fmt.Printf("Error creating uploads directory: %v\n", err)
		return
	}

	// Start listening on port
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server listening on port %s\n", PORT)
	fmt.Printf("Your IP address: %s\n", getLocalIP())

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Read the first line to determine the type of transfer
	firstLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading first line: %v\n", err)
		return
	}
	firstLine = strings.TrimSpace(firstLine)

	if strings.HasPrefix(firstLine, "FILE:") {
		// Handle file transfer
		handleFileTransfer(conn, reader, firstLine[5:])
	} else {
		// Handle text transfer
		handleTextTransfer(conn, reader)
	}
}

func handleFileTransfer(conn net.Conn, reader *bufio.Reader, filename string) {
	// Create the file in uploads directory
	filepath := filepath.Join("uploads", filename)
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	// Copy the file content
	_, err = io.Copy(file, reader)
	if err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		return
	}

	fmt.Printf("Received file: %s\n", filename)
}

func handleTextTransfer(conn net.Conn, reader *bufio.Reader) {
	// Read the text content
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading text: %v\n", err)
		return
	}

	fmt.Printf("Received text: %s\n", text)
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "unknown"
}
