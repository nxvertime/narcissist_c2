package structs

import "net"

type Channels struct {
	ServerIptCh chan bool
	ClientIptCh chan bool
}

type Client struct {
	ID      int
	Address string
	Conn    net.Conn
}
