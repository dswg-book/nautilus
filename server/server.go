package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
)

type ServerOptions struct {
	Redis      *RedisOptions
	Host, Port string
}

type Server struct {
	Options     *ServerOptions
	Context     context.Context
	listener    net.Listener
	connections map[net.Addr]*net.Conn
}

var serverInstance *Server

func NewServer(ctx context.Context, options *ServerOptions) *Server {
	serverInstance = &Server{
		Options:     options,
		Context:     ctx,
		connections: make(map[net.Addr]*net.Conn),
	}
	return serverInstance
}

func (s *Server) Start() {
	go func() {
		r := NewRedis(s.Context, *s.Options.Redis)
		if err := r.Connect(); err != nil {
			fmt.Printf("redis connect: error: %s\n", err)
			panic(err)
		}

		if err := s.listen(); err != nil {
			fmt.Printf("listen: error: %s\n", err)
			panic(err)
		}
	}()
}

func (s *Server) Stop() error {
	if err := redisInstance.Close(); err != nil {
		return err
	}
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
		s.connections[c.RemoteAddr()] = &c
		go s.handleConnection(c)
	}
}

func (s *Server) handleConnection(c net.Conn) {
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
			for _, conn := range s.connections {
				l := *conn
				l.Write([]byte(fmt.Sprintf("%s\n", input)))
			}
		}
	}
}
