package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/dswg-book/nautilus/config"
	"github.com/dswg-book/nautilus/peer"
	"github.com/dswg-book/nautilus/server"
)

type CLICommand struct {
	Command string
	Usage   string
	Args    []string
	Func    func(args ...string) error
}

func (c *CLICommand) Execute(args []string) error {
	return c.Func(args...)
}

func (c *CLICommand) String() string {
	return fmt.Sprintf("%s %s", c.Command, c.Usage)
}

var (
	currentServer *server.Server
	peers         = map[string]*peer.Peer{}
	currentConfig = config.NewConfig()
	commands      = map[string]*CLICommand{
		"add": {
			Command: "add",
			Usage:   "<peer_address>:<peer_port>",
			Args:    []string{},
			Func: func(args ...string) error {
				if len(args) != 1 {
					return fmt.Errorf("usage: add <peer_address>:<peer_port>")
				}
				peerAddress := args[0]
				peerAddressParts := strings.Split(peerAddress, ":")
				if len(peerAddressParts) != 2 {
					return fmt.Errorf("usage: add <peer_address>:<peer_port>")
				}
				peerAddress = peerAddressParts[0]
				peerPort, err := strconv.Atoi(peerAddressParts[1])
				if err != nil {
					return fmt.Errorf("usage: add <peer_address>:<peer_port>")
				}
				currentConfig.AddPeer(&peer.Peer{
					Host: peerAddress,
					Port: peerPort,
				})
				err = currentConfig.Save()
				if err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}
				return nil
			},
		},
		"remove": {
			Command: "remove",
			Usage:   "<peer_address>:<peer_port>",
			Args:    []string{},
			Func: func(args ...string) error {
				if len(args) != 1 {
					return fmt.Errorf("usage: remove <peer_address>:<peer_port>")
				}
				peerAddress := args[0]
				peerAddressParts := strings.Split(peerAddress, ":")
				if len(peerAddressParts) != 2 {
					return fmt.Errorf("usage: remove <peer_address>:<peer_port>")
				}
				peerHost := peerAddressParts[0]
				peerPort, err := strconv.Atoi(peerAddressParts[1])
				if err != nil {
					return fmt.Errorf("usage: remove <peer_address>:<peer_port>")
				}
				currentPeer := &peer.Peer{
					Host: peerHost,
					Port: peerPort,
				}
				currentConfig.RemovePeer(currentPeer)
				err = currentConfig.Save()
				if err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}
				return nil
			},
		},
		"list": {
			Command: "list",
			Usage:   "",
			Args:    []string{},
			Func: func(args ...string) error {
				err := currentConfig.Load()
				if err != nil {
					return fmt.Errorf("error loading config: %v", err)
				}
				for _, peer := range currentConfig.Peers {
					fmt.Printf("%s:%d\n", peer.Host, peer.Port)
				}
				return nil
			},
		},
		"start": {
			Command: "start",
			Usage:   "<host>:<port>",
			Args:    []string{},
			Func: func(args ...string) error {
				if len(args) != 1 {
					return fmt.Errorf("usage: start <host>:<port>")
				}
				address := args[0]
				addressParts := strings.Split(address, ":")
				if len(addressParts) != 2 {
					return fmt.Errorf("usage: start <host>:<port>")
				}
				host := addressParts[0]
				port, err := strconv.Atoi(addressParts[1])
				if err != nil {
					return fmt.Errorf("usage: start <host>:<port>")
				}

				signalChan := make(chan os.Signal, 1)
				signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

				currentServer = server.NewServer(
					context.Background(),
					&server.ServerOptions{
						Host: host,
						Port: port,
					},
				)
				for _, peer := range peers {
					peer.InputChannel = currentServer.InputChannel
					peer.OutputChannel = currentServer.OutputChannel
					peer.ErrorChannel = currentServer.ErrorChannel
					currentServer.AddPeer(peer)
				}

				currentServer.Start()
				fmt.Println("Server started on", address)

				<-signalChan
				currentServer.Stop()
				for _, peer := range peers {
					currentServer.RemovePeer(peer)
				}
				currentServer = nil
				return nil
			},
		},
	}
)

func main() {
	if len(os.Args) == 1 { // os.Args[0] is the program name
		fmt.Println("Usage: <command> [args]")
		return
	}

	command := os.Args[1]
	if commands[command] == nil {
		fmt.Println("Usage: <command> [args]")
		fmt.Println("Commands:")
		for command := range commands {
			fmt.Printf("  %-10s %s\n", command, commands[command].Usage)
		}
		fmt.Println("Example: start 0.0.0.0:3030")
		return
	}

	err := currentConfig.Load()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	for _, peer := range currentConfig.Peers {
		peers[fmt.Sprintf("%s:%d", peer.Host, peer.Port)] = peer
	}

	err = commands[command].Execute(os.Args[2:])
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
}
