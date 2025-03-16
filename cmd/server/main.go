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

		// Get and display the remote address
		remoteAddr := conn.RemoteAddr().String()
		fmt.Printf("New connection from: %s\n", remoteAddr)

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
		// Handle text transfer - the first line is the message
		fmt.Printf("Received text: %s\n", firstLine)
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

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}

	// First try to find LAN IPs (192.168.x.x, 10.x.x.x, 172.16.x.x)
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				// Check for common LAN IP patterns
				if ip4[0] == 192 && ip4[1] == 168 || // 192.168.x.x
					ip4[0] == 10 || // 10.x.x.x
					(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) { // 172.16.x.x to 172.31.x.x
					return ipnet.IP.String()
				}
			}
		}
	}

	// Fallback to any non-loopback IP
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				return ipnet.IP.String()
			}
		}
	}

	return "unknown"
}
