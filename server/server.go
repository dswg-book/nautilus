package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
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
			if !errors.Is(err, net.ErrClosed) {
				panic(err)
			}
		}
	}()
}

func (s *Server) Stop() error {
	return s.listener.Close()
}

func (s *Server) listen() error {
	l, err := net.Listen("tcp", s.Options.Host+":"+s.Options.Port)
	if err != nil {
		return err
	}
	s.listener = l

	log.Printf("server open at %s:%s\n", s.Options.Host, s.Options.Port)

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

func (s *Server) deleteConnection(c *Connection) {
	delete(s.connections, c.ID)
}

func (s *Server) closeAndDeleteConnection(c *Connection) {
	c.Close()
	s.deleteConnection(c)
}

func (s *Server) handleConnection(c *Connection) {
	if _, err := c.Write([]byte(fmt.Sprintf(">who:|>>id:%s\n", c.ID))); err != nil {
		log.Printf("connection write error: %s", err)
		s.closeAndDeleteConnection(c)
		return
	}

	for {
		data, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Printf("connection read error: %s", err)
			return
		}
		data = strings.TrimSpace(data)

		cmds := CommandsFromTags(data)
		for _, cmd := range cmds {
			if err := cmd.Run(c); err != nil {
				log.Printf("cmd run error: %s", err)
			}
		}
	}
}
