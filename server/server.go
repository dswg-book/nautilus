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
	Host string
	Port int
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
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Options.Host, s.Options.Port))
	if err != nil {
		return err
	}
	s.listener = l

	log.Printf("server open at %s:%d\n", s.Options.Host, s.Options.Port)

	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}
		conn := s.addConnection(c)
		log.Printf("[%s] connected", conn.ID)
		go s.handleConnection(conn)
	}
}

func (s *Server) hasID(id string) bool {
	var found bool
	for k, _ := range s.connections {
		if k == id {
			found = true
		}
	}
	return found
}
func (s *Server) addConnection(c net.Conn) *Connection {
	var name string
	var generatedName bool

	for !generatedName {
		name = generateName()
		generatedName = !s.hasID(name)
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

func (s *Server) updateConnection(c *Connection, cb func(*Connection)) {
	s.deleteConnection(c)
	cb(c)
	s.connections[c.ID] = c
}

func (s *Server) handleConnection(c *Connection) {
	if err := s.broadcast(fmt.Sprintf(">id:%s", c.ID)); err != nil {
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
		log.Printf("[%s] incoming data: %s", c.ID, data)

		cmds := CommandsFromTags(data)
		if len(cmds) < 1 {
			s.send(c, "", fmt.Sprintf(">message:%s: %s", CmdErrorInvalidCommand, data))
		}

		for _, cmd := range cmds {
			if err := cmd.Run(c); err != nil {
				log.Printf("cmd run error: %s", err)
			}
		}
	}
}

func (s *Server) broadcast(input string) error {
	for _, conn := range s.connections {
		if err := s.send(conn, "", input); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) send(conn *Connection, who string, input string) error {
	output := fmt.Sprintf(">who:%s|>%s\n", who, input)
	if _, err := conn.Write([]byte(output)); err != nil {
		log.Printf("conn: %s: send error: %s", conn.ID, err)
		return err
	}
	return nil
}
