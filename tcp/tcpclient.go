package tcp

import (
	"log"
	"fmt"
	"net"
	"errors"
	"github.com/nu7hatch/gouuid"
	"../session"
)

func ConnectToTcpServer(ip string, port int) (*net.TCPConn, error) {
	streamMap := make(map[*net.TCPConn]*Stream)
	streamChannel := make(chan *Stream)

	log.Println("Trying to connect to TCP server")

	tcpServerAddrString := fmt.Sprintf("%s:%d", ip, port)
	tcpServerAddr, err := net.ResolveTCPAddr("tcp", tcpServerAddrString)
	if err != nil { return nil, err }

	tcpConn, err := net.DialTCP("tcp", nil, tcpServerAddr)
	if err != nil { return nil, err }

	tcpConn.SetReadBuffer(0)
	tcpConn.SetWriteBuffer(0)

	go loopToReadStream(streamMap, streamChannel, tcpConn)
	go loopToHandleStream(streamMap, streamChannel)

	return tcpConn, nil
}

func Auth(tcpConn *net.TCPConn, sessionId *uuid.UUID) error {
	if tcpConn == nil {
		return errors.New("invalid tcpConn, tcpConn is nil")
	}

	currentSession, err := session.GetSession(sessionId)
	if err != nil {
		return err
	} else if currentSession == nil {
		return errors.New("invalid sessionId, couldn't find the session")
	}

	stream := append([]byte("AUTH"), (*sessionId)[:]...)
	_, err = tcpConn.Write(stream)
	if err != nil {
		return err
	}

	return nil
}
