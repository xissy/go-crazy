package client

import (
	"fmt"
	"net"
	"time"
	"errors"
	"../http"
	"../tcp"
	"../udp"
	"../session"
)

func ConnectToServer(host string, httpPort, tcpPort, udpPort int) (*session.Session, error){
	tcpConn, err := tcp.ConnectToTcpServer(host, tcpPort)
	if err != nil { return nil, err }

	baseUrl := fmt.Sprintf("http://%s:%d", host, httpPort)
	res, err := http.Connect(baseUrl)
	if err != nil { return nil, err }
	sessionId := res.SessionId

	currentSession := new(session.Session)
	currentSession.SessionId = sessionId
	session.PutSession(currentSession)

	err = tcp.Auth(tcpConn, sessionId)
	if err != nil { return nil, err }

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, udpPort))
	if err != nil { return nil, err }
	udpConn, err := udp.ConnectToUdpServer(host, udpPort)
	if err != nil { return nil, err }

	currentSession.TcpConn = tcpConn
	currentSession.UdpAddr = udpAddr
	currentSession.UdpConn = udpConn
	currentSession.HttpBaseUrl = baseUrl

	err = udp.SendPayloadsForInitialGap(currentSession)
	if err != nil { return nil, err }

	time.Sleep(100 * time.Millisecond)

	authResponseJson, err := http.Auth(baseUrl, sessionId)
	if err != nil { return nil, err }
	  
	fmt.Println("authResponseJson:", authResponseJson)

	if !authResponseJson.IsSuccess {
		return nil, errors.New("UDP payloads're not reached")
	}

	currentSession.SendingPayloadGap = authResponseJson.InitialPayloadGap
	currentSession.ReceivingPayloadGap = currentSession.InitialPayloadGap

	fmt.Println("currentSession:", currentSession)

	return currentSession, nil
}
