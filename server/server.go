package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
)

type ServerOptions struct {
	Host, Port string
}

type Server struct {
	Options     *ServerOptions
	Context     context.Context
	listener    net.Listener
	connections map[string]*Connection
}

var serverInstance *Server

func NewServer(ctx context.Context, options *ServerOptions) *Server {
	serverInstance = &Server{
		Options:     options,
		Context:     ctx,
		connections: make(map[string]*Connection),
	}
	return serverInstance
}

func (s *Server) Start() {
	go func() {
		if err := s.listen(); err != nil {
			fmt.Printf("listen: error: %s\n", err)
			panic(err)
		}
	}()
}

func (s *Server) Stop() error {
	return s.listener.Close()
}

func (s *Server) listen() error {
	l, err := net.Listen("tcp4", s.Options.Host+":"+s.Options.Port)
	if err != nil {
		return err
	}
	s.listener = l

	fmt.Printf("server open at %s:%s\n", s.Options.Host, s.Options.Port)

	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}
		conn := s.addConnection(c)
		go s.handleConnection(conn)
	}
}

func (s *Server) addConnection(c net.Conn) *Connection {
	var name string
	var generatedName bool

	for !generatedName {
		name = generateName()
		generatedName = true
		for k, _ := range s.connections {
			if k == name {
				generatedName = false
			}
		}
	}

	conn := NewConnection(ConnectionOptions{ID: name, Conn: c})
	s.connections[name] = conn
	return conn
}

func (s *Server) handleConnection(c *Connection) {
	if _, err := c.Write([]byte(fmt.Sprintf(">id:%s\n", c.ID))); err != nil {
		c.Close()
		return
	}

	for {
		data, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			return
		}
		data = strings.TrimSpace(data)

		input := data
		cmd := CmdMessage

		if strings.HasPrefix(data, "/") {
			parts := strings.SplitN(data, " ", 1)
			cut, _ := strings.CutPrefix(parts[0], "/")
			cmd = CmdCode(strings.ToLower(cut))
			input = parts[1]
		}

		switch cmd {
		case CmdDisconnect:
			c.Close()
			return
		case CmdMessage:
			for id, conn := range s.connections {
				l := *conn
				if id != c.ID {
					if _, err := l.Write([]byte(fmt.Sprintf(">who:%s >message:%s\n", c.ID, input))); err != nil {
						l.Close()
						return
					}
				}
			}
		}
	}
}
