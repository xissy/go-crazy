package udp

import (
	"fmt"
	"net"
	"log"
	"../payload"
)

func ConnectToUdpServer(host string, port int) (*net.UDPConn, error) {
	log.Println("Trying to connect to UDP server")

	udpServerAddrString := fmt.Sprintf("%s:%d", host, port)
	udpServerAddr, err := net.ResolveUDPAddr("udp", udpServerAddrString)
	if err != nil { return nil, err }

	udpConn, err := net.DialUDP("udp", nil, udpServerAddr)
	if err != nil { return nil, err }

	udpConn.SetReadBuffer(20000000)
	udpConn.SetWriteBuffer(20000000)

	payloadChannel := make(chan *payload.Payload)

	go loopToReadPayload(payloadChannel, udpConn)
	go loopToHandlePayload(payloadChannel)

	return udpConn, nil
}
