package server

import "net"

type ConnectionOptions struct {
	ID   string
	Conn net.Conn
}

type Connection struct {
	ID         string
	connection net.Conn
}

func NewConnection(options ConnectionOptions) *Connection {
	return &Connection{ID: options.ID, connection: options.Conn}
}

func (c *Connection) Read(p []byte) (bytesRead int, err error) {
	return c.connection.Read(p)
}

func (c *Connection) Write(p []byte) (bytesRead int, err error) {
	return c.connection.Write(p)
}

func (c *Connection) Close() error {
	return c.connection.Close()
}
