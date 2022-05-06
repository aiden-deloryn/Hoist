package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aiden-deloryn/hoist/src/client"
	"github.com/aiden-deloryn/hoist/src/server"
)

var (
	listenAddress = flag.String("listen", "localhost:8080", "port to listen to")
)

func main() {
	flag.Parse()

	go func() {
		time.Sleep(time.Second)

		if err := client.GetFileFromServer(*listenAddress); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	}()

	server.StartServer(*listenAddress)
}
