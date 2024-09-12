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
	host, port           string
	redisHost, redisPort string
)

func init() {
	flag.StringVar(&host, "host", "localhost", "Host for server")
	flag.StringVar(&port, "port", "3030", "Port for server")
	flag.StringVar(&redisHost, "redis-host", "127.0.0.1", "Host for Redis instance")
	flag.StringVar(&redisPort, "redis-port", "6379", "Port for Redis instance")
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
			Redis: &server.RedisOptions{
				Host: redisHost,
				Port: redisPort,
			},
		},
	)
	s.Start()

	<-signalChan
	s.Stop()
}
