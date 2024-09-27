package main

import (
	"bufio"
	"fmt"

	"github.com/DistributedSystemsWithGo/nautilus/client"
)

func main() {
	dialer := client.NewDialer("localhost", 3030)
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
		var cmd string
		if _, err := fmt.Scanln(&cmd); err != nil {
			fmt.Printf("error: %s\n", err.Error())
			continue
		}
		fmt.Fprintf(conn, "%s\n", cmd)
	}
}
