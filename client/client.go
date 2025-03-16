package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

const PORT = ":8080"

func SendText(serverIP, message string) {
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
	writer.Flush()

	fmt.Println("Message sent successfully")
}

func SendFile(serverIP, filePath string) {
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
