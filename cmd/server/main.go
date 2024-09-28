package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/DistributedSystemsWithGo/nautilus/server"
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
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	s := server.NewServer(
		context.Background(),
		&server.ServerOptions{
			Host: host,
			Port: port,
		},
	)
	s.Start()

	<-signalChan
	s.Stop()
}
