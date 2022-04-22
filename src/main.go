package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var (
	addr = flag.String("listen", "localhost:8080", "port to listen to")
	exit = false
)

func Server(address string) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	for !exit {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		HandleIncoming(conn)
	}
}

func HandleIncoming(conn io.ReadWriteCloser) {
	defer conn.Close()
	fmt.Println("Sending file...")
	SendFileToClient("file.txt", &conn)
	fmt.Println("File sent.")
	exit = true
}

func GetFileFromServer(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var filenameSize int64
	err = binary.Read(conn, binary.LittleEndian, &filenameSize)
	if err != nil {
		log.Println(err)
		return
	}

	filename := make([]byte, int(filenameSize))
	if _, err = io.ReadFull(conn, filename); err != nil {
		log.Println(err)
		return
	}

	var fileSize int64

	if err = binary.Read(conn, binary.LittleEndian, &fileSize); err != nil {
		log.Println(err)
		return
	}

	file, err := os.Create(string(filename) + ".copy")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	if err := file.Truncate(fileSize); err != nil {
		log.Println(err)
		return
	}

	br := bufio.NewReader(conn)
	bw := bufio.NewWriter(file)
	defer bw.Flush()

	if _, err = io.CopyN(bw, br, fileSize); err != nil {
		log.Println(err)
	}
}

func SendFileToClient(filename string, conn *io.ReadWriteCloser) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	filenameSize := int64(len(filename))

	if err = binary.Write(*conn, binary.LittleEndian, filenameSize); err != nil {
		log.Println(err)
		return
	}

	if _, err = io.WriteString(*conn, filename); err != nil {
		log.Println(err)
		return
	}

	stat, _ := file.Stat()
	if err = binary.Write(*conn, binary.LittleEndian, stat.Size()); err != nil {
		log.Println(err)
		return
	}

	br := bufio.NewReader(file)
	bw := bufio.NewWriter(*conn)
	defer bw.Flush()

	if _, err = io.CopyN(bw, br, stat.Size()); err != nil {
		log.Println(err)
		return
	}
}

func main() {
	flag.Parse()
	go func() {
		time.Sleep(time.Second)
		// SendFile(*addr, "file.txt")
		GetFileFromServer(*addr)
	}()
	Server(*addr)
}
