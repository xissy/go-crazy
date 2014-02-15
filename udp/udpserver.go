package udp

import (
	"fmt"
	"net"
	"log"
	"../payload"
)

func StartUdpServer(udpPort int) error {
	log.Println("Trying to start UDP server port:", udpPort)

	udpServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", udpPort))
	if err != nil { return err }

	udpConn, err := net.ListenUDP("udp", udpServerAddr)
	if err != nil { return err }

	udpConn.SetReadBuffer(20000000)
	udpConn.SetWriteBuffer(20000000)

	payloadChannel := make(chan *payload.Payload)

	go loopToReadPayload(payloadChannel, udpConn)
	go loopToHandlePayload(payloadChannel)

	return nil
}
