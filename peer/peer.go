package peer

import (
	"fmt"
	"net"
)

type Peer struct {
	Host          string       `json:"host"`
	Port          int          `json:"port"`
	connection    *net.Conn    `json:"-"`
	InputChannel  *chan string `json:"-"`
	OutputChannel *chan string `json:"-"`
	ErrorChannel  *chan error  `json:"-"`
	Connected     bool         `json:"-"`
}

func NewPeer(host string, port int, inputChannel *chan string, outputChannel *chan string, errorChannel *chan error) *Peer {
	peer := &Peer{
		Host:          host,
		Port:          port,
		InputChannel:  inputChannel,
		OutputChannel: outputChannel,
		ErrorChannel:  errorChannel,
		Connected:     false,
	}
	fmt.Printf("NewPeer: %+v\n", peer)
	return peer
}

func (p *Peer) Connect() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", p.Host, p.Port))
	if err != nil {
		*p.ErrorChannel <- err
	}
	p.connection = &conn
	p.Connected = true
	go func() {
		for message := range *p.InputChannel {
			p.Send(message)
		}
	}()
	go func() {
		for {
			if !p.Connected {
				break
			}
			p.Receive()
		}
	}()
}

func (p *Peer) Disconnect() {
	if p.connection == nil {
		return
	}
	err := (*p.connection).Close()
	if err != nil {
		*p.ErrorChannel <- err
	}
	p.connection = nil
	p.Connected = false
}

func (p *Peer) Send(message string) {
	if !p.Connected {
		*p.ErrorChannel <- fmt.Errorf("not connected to %s:%d", p.Host, p.Port)
		return
	}
	if p.connection == nil {
		*p.ErrorChannel <- fmt.Errorf("connection not established to %s:%d", p.Host, p.Port)
		return
	}
	_, err := (*p.connection).Write([]byte(message))
	if err != nil {
		*p.ErrorChannel <- err
	}
}

func (p *Peer) Receive() {
	if !p.Connected {
		*p.ErrorChannel <- fmt.Errorf("not connected to %s:%d", p.Host, p.Port)
		return
	}
	if p.connection == nil {
		*p.ErrorChannel <- fmt.Errorf("connection not established to %s:%d", p.Host, p.Port)
		return
	}
	buffer := make([]byte, 1024)
	n, err := (*p.connection).Read(buffer)
	if err != nil {
		*p.ErrorChannel <- err
		return
	}
	*p.OutputChannel <- string(buffer[:n])
}
