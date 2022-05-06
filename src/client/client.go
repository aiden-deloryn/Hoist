package client

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

func GetFileFromServer(address string) error {
	connection, err := net.Dial("tcp", address)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to connect to the server: %s", err))
	}

	defer connection.Close()

	// Receive the filename size from the server
	var filenameSize int64
	err = binary.Read(connection, binary.LittleEndian, &filenameSize)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read filename size from the server: %s", err))
	}

	// Receive the filename from the server
	filename := make([]byte, int(filenameSize))
	_, err = io.ReadFull(connection, filename)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read filename from the server: %s", err))
	}

	// Receive the file size from the server
	var fileSize int64
	err = binary.Read(connection, binary.LittleEndian, &fileSize)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read file size from the server: %s", err))
	}

	file, err := os.Create(string(filename) + ".copy")

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create file: %s", err))
	}

	defer file.Close()
	err = file.Truncate(fileSize)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to set file size: %s", err))
	}

	reader := bufio.NewReader(connection)
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Receive the file from the server
	_, err = io.CopyN(writer, reader, fileSize)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to receive file from the server: %s", err))
	}

	return nil
}
