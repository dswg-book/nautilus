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

func (s *Server) deleteConnection(c *Connection) {
	delete(s.connections, c.ID)
}

func (s *Server) closeAndDeleteConnection(c *Connection) {
	c.Close()
	s.deleteConnection(c)
}

func (s *Server) handleConnection(c *Connection) {
	if _, err := c.Write([]byte(fmt.Sprintf(">who:|>>id:%s\n", c.ID))); err != nil {
		c.Close()
		return
	}

	for {
		data, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			return
		}
		data = strings.TrimSpace(data)

		var cmds []*Command

		tags := strings.Split(data, "|>")
		for _, tag := range tags {
			cmd := NewCommand(CommandOptions{Code: CmdMessage})
			if strings.HasPrefix(tag, "<") {
				tagParts := strings.Split(tag, ":")
				cmd = NewCommand(
					CommandOptions{
						Code:  CmdCode(tagParts[0][1:]),
						Input: tagParts[1],
					},
				)
			}
			cmds = append(cmds, cmd)
		}

		for _, cmd := range cmds {
			if cmd.Code == CmdAction {
				actionParts := strings.SplitN(cmd.Input, " ", 1)
				cmd.Input = ""
				cmd.Code = CmdCode(actionParts[0])
				if len(actionParts) > 1 {
					cmd.Input = actionParts[1]
				}
			}

			switch cmd.Code {
			case CmdDisconnect, CmdClose:
				s.closeAndDeleteConnection(c)
				return
			case CmdMessage:
				for id, conn := range s.connections {
					l := *conn
					if id != c.ID {
						if cmd.Input != "" {
							if _, err := l.Write([]byte(fmt.Sprintf(">who:%s|>>message:%s\n", c.ID, cmd.Input))); err != nil {
								serverInstance.closeAndDeleteConnection(conn)
								return
							}
						}
					}
				}
			default:
				if _, err := c.Write([]byte(fmt.Sprintf(">who:|>>message:%s\n", CmdErrorInvalidCommand))); err != nil {
					s.closeAndDeleteConnection(c)
					return
				}
			}
		}
	}
}
