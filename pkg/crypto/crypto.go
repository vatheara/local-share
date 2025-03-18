package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"syscall"

	"golang.org/x/term"
)

// GetEncryptionKey retrieves the encryption key from environment variable or prompts user
func GetEncryptionKey(confirmPassword bool) (string, error) {
	// First try environment variable
	key := os.Getenv("LOCALSHARE_KEY")
	if key != "" {
		return PadKey(key), nil
	}

	// If no environment variable, prompt user
	if confirmPassword {
		fmt.Println("Please enter a password to encrypt messages:")
	} else {
		fmt.Println("Please enter the password to decrypt messages:")
	}
	fmt.Print("Password: ")
	keyBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Add newline after password input
	if err != nil {
		return "", fmt.Errorf("error reading password: %v", err)
	}

	if confirmPassword {
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
	}

	return PadKey(string(keyBytes)), nil
}

// PadKey ensures the key is exactly 32 bytes by padding or truncating
func PadKey(key string) string {
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

// Encrypt encrypts plaintext with AES-256
func Encrypt(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Generate a random IV
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Encrypt
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// Return as base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext with AES-256
func Decrypt(encryptedMsg string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedMsg)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Check ciphertext length
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Extract IV
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Decrypt
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
} 