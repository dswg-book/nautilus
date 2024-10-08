package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/DistributedSystemsWithGo/nautilus/client"
)

var (
	host string
	port int
)

func init() {
	flag.StringVar(&host, "host", "localhost", "Host for server")
	flag.IntVar(&port, "port", 3030, "Port for server")
	flag.Parse()
}

func main() {
	dialer := client.NewDialer(host, port)
	conn, err := dialer.Open()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			incoming, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				panic(err)
			}
			fmt.Print(incoming)
		}
	}()
	for {
		in := bufio.NewReader(os.Stdin)
		cmd, err := in.ReadString('\n')
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			continue
		}
		fmt.Fprintf(conn, "%s\n", cmd)
	}
}
