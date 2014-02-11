package tcp

import (
	"log"
	"fmt"
	"net"
)

func StartTcpServer(tcpPort int) error {
	log.Println("Trying to start TCP server port:", tcpPort)
	
	tcpServerAddrString := fmt.Sprintf("0.0.0.0:%d", tcpPort)
	tcpServerAddr, err := net.ResolveTCPAddr("tcp", tcpServerAddrString)
	if err != nil { return err }

	tcpListener, err := net.ListenTCP("tcp", tcpServerAddr)
	if err != nil { return err }

	streamMap := make(map[*net.TCPConn]*Stream)
	streamChannel := make(chan *Stream)

	go loopToAccept(streamMap, streamChannel, tcpListener)
	go loopToHandleStream(streamMap, streamChannel)

	return nil
}

func loopToAccept(streamMap map[*net.TCPConn]*Stream,
				streamChannel chan *Stream,
				tcpListener *net.TCPListener) error {
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil { return err }

		tcpConn.SetReadBuffer(0)
		tcpConn.SetWriteBuffer(0)

		go loopToReadStream(streamMap, streamChannel, tcpConn)
	}

	return nil
}
