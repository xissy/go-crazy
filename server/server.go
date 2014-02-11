package server

import (
	"../tcp"
	"../http"
	"../udp"
)

func StartServer(httpPort, tcpPort, udpPort int) error {
	err := tcp.StartTcpServer(tcpPort)
	if err != nil { return err }

	err = http.StartHttpServer(httpPort)
	if err != nil { return err }

	err = udp.StartUdpServer(udpPort)
	if err != nil { return err }

	return nil
}
