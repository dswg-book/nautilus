package client

import (
	"fmt"
	"net"
)

type Dialer struct {
	host string
	port int
}

func NewDialer(host string, port int) *Dialer {
	return &Dialer{host, port}
}

func (d *Dialer) Open() (net.Conn, error) {
	return net.Dial("tcp", d.address())
}

func (d *Dialer) address() string {
	return fmt.Sprintf("%s:%d", d.host, d.port)
}
