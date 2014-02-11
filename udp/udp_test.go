package udp

import (
	"testing"
	"net"
	"time"
	"github.com/nu7hatch/gouuid"
)

var _udpConn *net.UDPConn

func TestStartUdpServer(t *testing.T) {
	udpPort := 20069

	err := StartUdpServer(udpPort)
	if err != nil {
		t.Error("Failed to Start UDP Server")
	}
}

func TestConnectToUdpServer(t *testing.T) {
	udpIp := "127.0.0.1"
	udpPort := 20069

	udpConn, err := ConnectToUdpServer(udpIp, udpPort)
	if err != nil {
		t.Error("Failed to Connect")
	}

	_udpConn = udpConn
}

func TestWriteToServer(t *testing.T) {
	sessionId, err := uuid.NewV4()
	if err != nil {
		t.Error("Failed to Create an UUID for SessionId")
	}

	payload := append([]byte("4GAP"), (*sessionId)[:]...)
	_udpConn.Write(payload)

	time.Sleep(100 * time.Millisecond)
}
