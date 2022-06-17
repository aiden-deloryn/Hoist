package client

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/aiden-deloryn/hoist/src/util"
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

	file, err := os.Create(string(filename))

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create file: %s", err))
	}

	defer file.Close()
	err = file.Truncate(fileSize)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to set file size: %s", err))
	}

	reader := bufio.NewReader(connection)

	// Init vars to measure copy speed
	copySpeed := int64(0)
	sampleStartTime := time.Now().UnixMilli() - 1
	sampleStartBytes := int64(0)

	// Convert our bufio.Reader into a util.ProgressReader so we can log the
	// progress of a copy to the console.
	progressReader := &util.ProgressReader{
		Reader: *reader,
		ProgressCallback: func(bytesCopied int64) {
			copyComplete := bytesCopied == fileSize
			progress := int(float64(bytesCopied) / float64(fileSize) * 100)

			sampleDuration := time.Now().UnixMilli() - sampleStartTime
			sampleBytesCopied := bytesCopied - sampleStartBytes

			// Calculate the current copy speed
			if sampleDuration >= 1000 {
				// Calculate speed in MiB per second
				copySpeed = (sampleBytesCopied / 1048576) / (sampleDuration / 1000)
				sampleStartTime = time.Now().UnixMilli()
				sampleStartBytes = bytesCopied
			}

			fmt.Printf("\rCopying file %s %d/%d bytes (%d MiB/s)", util.GenerateProgressBarString(progress), bytesCopied, fileSize, copySpeed)

			if copyComplete {
				fmt.Print("\n")
			}
		},
	}

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Receive the file from the server
	_, err = io.CopyN(writer, progressReader, fileSize)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to receive file from the server: %s", err))
	}

	return nil
}
