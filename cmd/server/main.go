package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/term"
)

const (
	PORT        = ":8080"
	BUFFER_SIZE = 1024 * 1024 // 1MB buffer for file transfers
)

// getEncryptionKey retrieves the encryption key from environment variable or prompts user
func getEncryptionKey() (string, error) {
	// First try environment variable
	key := os.Getenv("LOCALSHARE_KEY")
	if key != "" {
		return padKey(key), nil
	}

	// If no environment variable, prompt user
	fmt.Println("Please enter a password to encrypt messages:")
	fmt.Print("Password: ")
	keyBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Add newline after password input
	if err != nil {
		return "", fmt.Errorf("error reading password: %v", err)
	}

	// Confirm password
	fmt.Print("Confirm password: ")
	confirmBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Add newline after password input
	if err != nil {
		return "", fmt.Errorf("error reading password confirmation: %v", err)
	}

	if string(keyBytes) != string(confirmBytes) {
		return "", fmt.Errorf("passwords do not match")
	}

	return padKey(string(keyBytes)), nil
}

// padKey ensures the key is exactly 32 bytes by padding or truncating
func padKey(key string) string {
	if len(key) == 0 {
		// If empty, use a default key
		return "default-32-byte-key-for-local-share!!"
	}

	if len(key) >= 32 {
		// If longer than 32 bytes, truncate
		return key[:32]
	}

	// If shorter than 32 bytes, pad with the key itself
	padded := make([]byte, 32)
	copy(padded, key)
	for i := len(key); i < 32; i++ {
		padded[i] = padded[i-len(key)]
	}
	return string(padded)
}

func main() {
	// Get the encryption key
	encryptionKey, err := getEncryptionKey()
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
		// Handle file transfer
		handleFileTransfer(conn, reader, firstLine[5:])
	} else if strings.HasPrefix(firstLine, "TEXT:") {
		// Handle encrypted text transfer
		encryptedMsg := firstLine[5:]
		decryptedMsg, err := decrypt(encryptedMsg, []byte(encryptionKey))
		if err != nil {
			fmt.Printf("Error decrypting message: %v\n", err)
			return
		}
		fmt.Printf("Received decrypted text: %s\n", decryptedMsg)
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

// Decryption helper function
func decrypt(encryptedMsg string, key []byte) (string, error) {
	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedMsg)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract nonce size
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the message
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
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
