package sender

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"local-share/pkg/crypto"
)

const (
	PORT = ":8080"
)

// SendText sends encrypted text to a server
func SendText(serverIP, message string) {
	// Get the encryption key
	key, err := crypto.GetEncryptionKey(false)
	if err != nil {
		fmt.Printf("Error getting encryption key: %v\n", err)
		return
	}

	conn, err := net.Dial("tcp", serverIP+PORT)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer conn.Close()

	// Encrypt the message
	encryptedMsg, err := crypto.Encrypt([]byte(message), []byte(key))
	if err != nil {
		fmt.Printf("Error encrypting message: %v\n", err)
		return
	}

	// Send the encrypted message with "TEXT:" prefix
	writer := bufio.NewWriter(conn)
	_, err = writer.WriteString("TEXT:" + encryptedMsg + "\n")
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return
	}

	// Ensure the message is sent
	if err := writer.Flush(); err != nil {
		fmt.Printf("Error flushing message: %v\n", err)
		return
	}

	fmt.Println("Encrypted message sent successfully")
}

// SendFile sends an encrypted file to a server
func SendFile(serverIP, filePath string) {
	// Get the encryption key
	key, err := crypto.GetEncryptionKey(false)
	if err != nil {
		fmt.Printf("Error getting encryption key: %v\n", err)
		return
	}

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

	// Read the entire file
	fileContent, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Encrypt the file content
	encryptedContent, err := crypto.Encrypt(fileContent, []byte(key))
	if err != nil {
		fmt.Printf("Error encrypting file: %v\n", err)
		return
	}

	// Encrypt the filename
	filename := filepath.Base(filePath)
	encryptedFilename, err := crypto.Encrypt([]byte(filename), []byte(key))
	if err != nil {
		fmt.Printf("Error encrypting filename: %v\n", err)
		return
	}

	writer := bufio.NewWriter(conn)

	// Send the encrypted filename
	_, err = writer.WriteString("FILE:" + encryptedFilename + "\n")
	if err != nil {
		fmt.Printf("Error sending filename: %v\n", err)
		return
	}
	writer.Flush()

	// Send the encrypted content length followed by content
	contentLength := len(encryptedContent)
	_, err = writer.WriteString(fmt.Sprintf("%d\n", contentLength))
	if err != nil {
		fmt.Printf("Error sending content length: %v\n", err)
		return
	}
	writer.Flush()

	// Send the encrypted content
	_, err = writer.WriteString(encryptedContent)
	if err != nil {
		fmt.Printf("Error sending file content: %v\n", err)
		return
	}
	writer.Flush()

	fmt.Printf("File %s encrypted and sent successfully\n", filename)
} 