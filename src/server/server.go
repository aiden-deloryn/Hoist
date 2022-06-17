package server

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

func StartServer(address string, filename string, keepAlive bool) {
	listner, err := net.Listen("tcp", address)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for {
		conn, err := listner.Accept()

		// If a connection error occurs, log the error and move on to the next connection
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if keepAlive {
			// Handle multiple connections by starting a new goroutine for each one
			go HandleIncomingConnection(conn, filename)
		} else {
			// Handle the first successful connection and then exit
			HandleIncomingConnection(conn, filename)
			break
		}
	}
}

func HandleIncomingConnection(connection net.Conn, filename string) {
	defer connection.Close()
	fmt.Printf("Sending file to %s...\n", connection.RemoteAddr())
	err := SendFileToClient(filename, &connection)

	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("An error occurred when sending file: %s", err))
		return
	}

	fmt.Printf("File sent to %s\n", connection.RemoteAddr())
}

func SendFileToClient(filename string, conn *net.Conn) error {
	file, err := os.Open(filename)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read file: %s", err))
	}

	defer file.Close()

	destFilename := filepath.Base(filename)

	filenameSize := int64(len(destFilename))

	// Send the size of the filename to the client
	err = binary.Write(*conn, binary.LittleEndian, filenameSize)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to send filename size to the client: %s", err))
	}

	// Send the filename to the client
	_, err = io.WriteString(*conn, destFilename)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to send filename to the client: %s", err))
	}

	fileInfo, err := file.Stat()

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get file info: %s", err))
	}

	// Send the size of the file to the client
	err = binary.Write(*conn, binary.LittleEndian, fileInfo.Size())

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to send file size to the client: %s", err))
	}

	reader := bufio.NewReader(file)
	writer := bufio.NewWriter(*conn)
	defer writer.Flush()

	// Send the file to the client
	_, err = io.CopyN(writer, reader, fileInfo.Size())

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to send file to the client: %s", err))
	}

	return nil
}
