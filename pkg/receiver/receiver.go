package receiver

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"local-share/pkg/crypto"
)

const (
	PORT        = ":8080"
	BUFFER_SIZE = 1024 * 1024 // 1MB buffer for file transfers
)

// Start starts the receiver server
func Start() {
	// Get the encryption key
	encryptionKey, err := crypto.GetEncryptionKey(true)
	if err != nil {
		fmt.Printf("Error getting encryption key: %v\n", err)
		return
	}

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

		go handleConnection(conn, encryptionKey)
	}
}

func handleConnection(conn net.Conn, encryptionKey string) {
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
		// Handle encrypted file transfer
		handleFileTransfer(conn, reader, firstLine[5:], encryptionKey)
	} else if strings.HasPrefix(firstLine, "TEXT:") {
		// Handle encrypted text transfer
		encryptedMsg := firstLine[5:]
		decryptedMsg, err := crypto.Decrypt(encryptedMsg, []byte(encryptionKey))
		if err != nil {
			fmt.Printf("Error decrypting message: %v\n", err)
			return
		}
		fmt.Printf("Received decrypted text: %s\n", decryptedMsg)
	}
}

func handleFileTransfer(conn net.Conn, reader *bufio.Reader, encryptedFilename string, encryptionKey string) {
	// Decrypt the filename
	filename, err := crypto.Decrypt(encryptedFilename, []byte(encryptionKey))
	if err != nil {
		fmt.Printf("Error decrypting filename: %v\n", err)
		return
	}

	// Read the content length
	lengthStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading content length: %v\n", err)
		return
	}
	contentLength := 0
	_, err = fmt.Sscanf(strings.TrimSpace(lengthStr), "%d", &contentLength)
	if err != nil {
		fmt.Printf("Error parsing content length: %v\n", err)
		return
	}

	// Read the encrypted content
	encryptedContent := make([]byte, contentLength)
	_, err = io.ReadFull(reader, encryptedContent)
	if err != nil {
		fmt.Printf("Error reading file content: %v\n", err)
		return
	}

	// Decrypt the content
	decryptedContent, err := crypto.Decrypt(string(encryptedContent), []byte(encryptionKey))
	if err != nil {
		fmt.Printf("Error decrypting file content: %v\n", err)
		return
	}

	// Create the file in uploads directory
	filepath := filepath.Join("uploads", filename)
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	// Write the decrypted content
	_, err = file.Write([]byte(decryptedContent))
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("Received and decrypted file: %s\n", filename)
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
					return ip4.String()
				}
			}
		}
	}

	// If no LAN IP found, return any non-loopback IP
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				return ip4.String()
			}
		}
	}

	return "127.0.0.1" // Loopback as fallback
}
