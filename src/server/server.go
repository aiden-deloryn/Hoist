package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/aiden-deloryn/hoist/src/values"
)

func StartServer(address string, filename string, password string, keepAlive bool) error {
	listner, err := net.Listen("tcp", address)

	if err != nil {
		return fmt.Errorf("failed to start TCP server: %s", err)
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
			go handleIncomingConnection(conn, filename, password)
		} else {
			// Handle the first successful connection and then exit
			handleIncomingConnection(conn, filename, password)
			break
		}
	}

	return nil
}

func handleIncomingConnection(conn net.Conn, filename string, password string) error {
	defer conn.Close()

	err := verifyPassword(conn, password)

	if err != nil {
		return fmt.Errorf("Failed to verify password: %s", err)
	}

	fmt.Printf("Sending file(s) to %s...\n", conn.RemoteAddr())
	err = sendObjectToClient(filename, conn)

	if err != nil {
		return fmt.Errorf("An error occurred when sending file: %s", err)
	}

	fmt.Printf("File(s) sent to %s\n", conn.RemoteAddr())

	return nil
}

func verifyPassword(conn net.Conn, password string) error {
	// Pad the password with '0' so that it's length is MAX_PASSWORD_LENGTH
	for len(password) < values.MAX_PASSWORD_LENGTH {
		password += "0"
	}

	// Get the password from the client
	guess := bytes.NewBuffer([]byte{})
	_, err := io.CopyN(guess, conn, int64(len(password)))

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read password from the client: %s", err))
	}

	if string(string(guess.Bytes())) != password {
		// Notify the client that password verification failed
		io.CopyN(conn, bytes.NewBuffer([]byte{0}), 1)
		return fmt.Errorf("Password is incorrect")
	}

	// Notify the client that password verification succeeded
	_, err = io.CopyN(conn, bytes.NewBuffer([]byte{1}), 1)

	if err != nil {
		return fmt.Errorf("Failed to notify client of password verification result: %s", err)
	}

	return nil
}

func sendObjectToClient(filename string, conn net.Conn) error {
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

			destFilename := strings.TrimPrefix(path, filepath.Dir(filename))
			destFilename = strings.TrimPrefix(destFilename, string(filepath.Separator))

			err = sendFileToClient(path, destFilename, conn)

			if err != nil {
				return errors.New(fmt.Sprintf("Failed to send file to client '%s': %s", path, err))
			}

			return nil
		})
	} else {
		// Send file
		destFilename := filepath.Base(filename)
		err = sendFileToClient(filename, destFilename, conn)
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

func sendFileToClient(srcFilename string, destFilename string, conn net.Conn) error {
	file, err := os.Open(srcFilename)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read file: %s", err))
	}

	defer file.Close()

	// Always send paths using the '/' separator over the network. These paths
	// will be converted to platform specific paths by the client
	destFilename = strings.ReplaceAll(destFilename, string(filepath.Separator), "/")

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
