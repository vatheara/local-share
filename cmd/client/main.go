package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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
	PORT = ":8080"
)

// getEncryptionKey retrieves the encryption key from environment variable or prompts user
func getEncryptionKey() (string, error) {
	// First try environment variable
	key := os.Getenv("LOCALSHARE_KEY")
	if key != "" {
		return padKey(key), nil
	}

	// If no environment variable, prompt user
	fmt.Println("Please enter the password to decrypt messages:")
	fmt.Print("Password: ")
	keyBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Add newline after password input
	if err != nil {
		return "", fmt.Errorf("error reading password: %v", err)
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
	// Get the encryption key
	key, err := getEncryptionKey()
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
	encryptedMsg, err := encrypt([]byte(message), []byte(key))
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

// Encryption helper function
func encrypt(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt and seal the data
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Convert to base64 for safe transmission
	return base64.StdEncoding.EncodeToString(ciphertext), nil
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
