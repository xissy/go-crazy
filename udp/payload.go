package udp

import (
	"net"
)

type Payload struct {
	UdpAddr *net.UDPAddr
	UdpConn *net.UDPConn
	Buffer []byte
	BufferLength int
	Err error
}
