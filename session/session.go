package session

import (
	"net"
	"time"
	"github.com/nu7hatch/gouuid"
)

type Session struct {
	SessionId *uuid.UUID
	UdpAddr *net.UDPAddr
	UdpConn *net.UDPConn
	TcpConn *net.TCPConn
	HttpBaseUrl string
	IsDisconnected bool
	InitialPayloadGap time.Duration
	InitialPayloadGapSum time.Duration
	InitialPayloadCount int
	PrevInitialPayloadTime time.Time
	SendingPayloadGap time.Duration
	ReceivingPayloadGap time.Duration
}
