package tcp

import (
	"net"
	"github.com/nu7hatch/gouuid"
)

type Stream struct {
	TcpConn *net.TCPConn
	Buffer []byte
	BufferLength int
	SessionId *uuid.UUID
}
