package tcp

import (
	"testing"
	"net"
	"time"
	"github.com/nu7hatch/gouuid"
)

var _tcpConn *net.TCPConn

func TestStartTcpServer(t *testing.T) {
	tcpPort := 20068

	err := StartTcpServer(tcpPort)
	if err != nil {
		t.Error("Failed to Start TCP Server")
	}
}

func TestConnectToTcpServer(t *testing.T) {
	tcpIp := "127.0.0.1"
	tcpPort := 20068

	tcpConn, err := ConnectToTcpServer(tcpIp, tcpPort)
	if err != nil {
		t.Error("Failed to Connect")
	}

	_tcpConn = tcpConn
}

func TestWriteToServer(t *testing.T) {
	_tcpConn.Write( append([]byte("AUTH"), (*new(uuid.UUID))[:]...) )

	time.Sleep(1 * time.Second)
}
