package payload

import (
	"net"
)

const DefaultPayloadDataSize int = 1024

type Payload struct {
	UdpAddr *net.UDPAddr
	UdpConn *net.UDPConn
	Buffer []byte
	BufferLength int
	Err error
}
