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
	"strings"
)

func StartServer(address string, filename string, keepAlive bool) {
	listner, err := net.Listen("tcp", address)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Printf("The target file or directory is ready to send. To download it on another machine, use:\n")
	fmt.Printf("  hoist get %s\n", address)

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
	fmt.Printf("Sending file(s) to %s...\n", connection.RemoteAddr())
	err := SendObjectToClient(filename, connection)

	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("An error occurred when sending file: %s", err))
		return
	}

	fmt.Printf("File(s) sent to %s\n", connection.RemoteAddr())
}

func SendObjectToClient(filename string, conn net.Conn) error {
	file, err := os.Open(filename)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read file: %s", err))
	}

	defer file.Close()

	fileInfo, err := file.Stat()

	if fileInfo.IsDir() {
		// Send directory
		err = filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			destFilename := strings.TrimPrefix(path, filepath.Dir(filename)+string(filepath.Separator))

			err = SendFileToClient(path, destFilename, conn)

			if err != nil {
				return errors.New(fmt.Sprintf("Failed to send file to client '%s': %s", path, err))
			}

			return nil
		})
	} else {
		// Send file
		destFilename := filepath.Base(filename)
		err = SendFileToClient(filename, destFilename, conn)
	}

	if err != nil {
		return err
	}

	// Notify the client there is nothing left to copy by sending a file size of -1
	err = binary.Write(conn, binary.LittleEndian, int64(-1))

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to send 'copy complete' message to client: %s", err))
	}

	return nil
}

func SendFileToClient(srcFilename string, destFilename string, conn net.Conn) error {
	file, err := os.Open(srcFilename)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read file: %s", err))
	}

	defer file.Close()

	filenameSize := int64(len(destFilename))

	// Send the size of the filename to the client
	err = binary.Write(conn, binary.LittleEndian, filenameSize)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to send filename size to the client: %s", err))
	}

	// Send the filename to the client
	_, err = io.WriteString(conn, destFilename)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to send filename to the client: %s", err))
	}

	fileInfo, err := file.Stat()

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get file info: %s", err))
	}

	// Send the size of the file to the client
	err = binary.Write(conn, binary.LittleEndian, fileInfo.Size())

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to send file size to the client: %s", err))
	}

	reader := bufio.NewReader(file)

	// Send the file to the client
	_, err = io.CopyN(conn, reader, fileInfo.Size())

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to send file to the client: %s", err))
	}

	return nil
}
