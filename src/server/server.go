package server

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

var (
	keepAlive = false
)

func StartServer(address string, filename string) {
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

func HandleIncomingConnection(connection io.ReadWriteCloser, filename string) {
	defer connection.Close()
	fmt.Println("Sending file...")
	err := SendFileToClient(filename, &connection)

	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("An error occurred when sending file: %s", err))
		return
	}

	fmt.Println("File sent.")
}

func SendFileToClient(filename string, conn *io.ReadWriteCloser) error {
	file, err := os.Open(filename)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read file: %s", err))
	}

	defer file.Close()

	filenameSize := int64(len(filename))

	// Send the size of the filename to the client
	err = binary.Write(*conn, binary.LittleEndian, filenameSize)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to send filename size to the client: %s", err))
	}

	// Send the filename to the client
	_, err = io.WriteString(*conn, filename)

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
