package client

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aiden-deloryn/hoist/src/types"
	"github.com/aiden-deloryn/hoist/src/util"
	"github.com/aiden-deloryn/hoist/src/values"
)

func GetFileFromServer(address string, password string, outputDirectory string) error {
	conn, err := net.Dial("tcp", address)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to connect to the server: %s", err))
	}

	defer conn.Close()

	err = authenticate(conn, password)

	if err != nil {
		return fmt.Errorf("Authentication failed: %s", err)
	}

	for {
		// Receive the filename size from the server
		var filenameSize int64
		err = binary.Read(conn, binary.LittleEndian, &filenameSize)

		if err != nil {
			return errors.New(fmt.Sprintf("Failed to read filename size from the server: %s", err))
		}

		// If the filenameSize is -1, there is nothing left to copy
		if filenameSize == -1 {
			break
		}

		// If the filenameSize is -2, the next object is a symlink
		if filenameSize == -2 {
			err = GetSymlinkFromServer(conn, outputDirectory)

			if err != nil {
				return fmt.Errorf("failed to get symlink from server: %s", err)
			}

			continue
		}

		// Receive the filenameBytes from the server
		filenameBytes := make([]byte, int(filenameSize))
		_, err = io.ReadFull(conn, filenameBytes)

		if err != nil {
			return errors.New(fmt.Sprintf("Failed to read filename from the server: %s", err))
		}

		// Convert filename's path separator for the current platform
		filename := filepath.FromSlash(string(filenameBytes))

		if outputDirectory != "" {
			filename = filepath.Clean(outputDirectory + string(filepath.Separator) + filename)
		}

		// Create parent directories
		if strings.Count(filename, string(filepath.Separator)) != 0 {
			os.MkdirAll(filepath.Dir(filename), 0775)
		}

		// Receive the file size from the server
		var fileSize int64
		err = binary.Read(conn, binary.LittleEndian, &fileSize)

		if err != nil {
			return errors.New(fmt.Sprintf("Failed to read file size from the server: %s", err))
		}

		file, err := os.Create(filename)

		if err != nil {
			return errors.New(fmt.Sprintf("Failed to create file: %s", err))
		}

		defer file.Close()
		err = file.Truncate(fileSize)

		if err != nil {
			return errors.New(fmt.Sprintf("Failed to set file size: %s", err))
		}

		// Init vars to measure copy speed
		copySpeed := int64(0)
		sampleStartTime := time.Now().UnixMilli() - 1
		sampleStartBytes := int64(0)

		// Convert our bufio.Reader into a util.ProgressReader so we can log the
		// progress of a copy to the console.
		progressReader := &util.ProgressReader{
			Reader: conn,
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

				fmt.Printf("\r%s %d/%d bytes (%d MiB/s)", util.GenerateProgressBarString(progress), bytesCopied, fileSize, copySpeed)

				if copyComplete {
					fmt.Print("\n")
				}
			},
		}

		writer := bufio.NewWriter(file)

		fmt.Printf("Copying file %s...\n", filename)

		// Receive the file from the server
		_, err = io.CopyN(writer, progressReader, fileSize)

		if err != nil {
			return errors.New(fmt.Sprintf("Failed to receive file from the server: %s", err))
		}

		file.Close()
	}

	return nil
}

func GetSymlinkFromServer(conn net.Conn, outputDirectory string) error {
	// Receive the symlink metadata size from the server
	var metadataSize int64
	err := binary.Read(conn, binary.LittleEndian, &metadataSize)

	if err != nil {
		return fmt.Errorf("failed to read symlink metadata size from the server: %s", err)
	}

	// Receive the symlink metadata from the server
	symlinkMetadataJSON := make([]byte, int(metadataSize))
	_, err = io.ReadFull(conn, symlinkMetadataJSON)

	if err != nil {
		return fmt.Errorf("failed to read symlink metadata from the server: %s", err)
	}

	metadata := types.SymlinkMetadata{}

	err = json.Unmarshal(symlinkMetadataJSON, &metadata)

	if err != nil {
		return fmt.Errorf("failed to unmarshal symlink metadata: %s", err)
	}

	// Convert filename's path separator for the current platform
	metadata.Target = filepath.FromSlash(metadata.Target)
	metadata.Name = filepath.FromSlash(metadata.Name)

	if outputDirectory != "" {
		metadata.Name = filepath.Clean(outputDirectory + string(filepath.Separator) + metadata.Name)
	}

	os.MkdirAll(filepath.Dir(metadata.Name), 0775)

	fmt.Printf("Creating symlink: \n")
	fmt.Printf("  %s --> %s\n", metadata.Name, metadata.Target)

	err = os.Symlink(metadata.Target, metadata.Name)

	if err != nil {
		return fmt.Errorf("failed to create symlink: %s", err)
	}

	return nil
}

func authenticate(conn net.Conn, password string) error {
	// Pad the password with '0' so that it's length is MAX_PASSWORD_LENGTH
	for len(password) < values.MAX_PASSWORD_LENGTH {
		password += "0"
	}

	// Send the password to the server
	_, err := io.CopyN(conn, bytes.NewBuffer([]byte(password)), int64(len(password)))

	if err != nil {
		return fmt.Errorf("Failed to send data to server: %s", err)
	}

	// Read the result from the server (result is boolean 0 or 1)
	result := bytes.NewBuffer([]byte{})
	_, err = io.CopyN(result, conn, 1)

	if err != nil {
		return fmt.Errorf("Failed to get response from server: %s", err)
	}

	// If the result was 0 (false) we can assume the password was incorrect
	if result.Bytes()[0] == 0 {
		return fmt.Errorf("Password is incorrect")
	}

	return nil
}
