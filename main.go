package main

import (
	"fmt"
	"time"
	"./server"
	"./file"
)

func main() {
	httpPort := 8080
	tcpPort := 20068
	udpPort := 20069

	file.SendingFileMap = make(file.FileMap)
	file.ReceivingFileMap = make(file.FileMap)

	err := server.StartServer(httpPort, tcpPort, udpPort)
	if err != nil {
	  fmt.Println("Failed to Start the Server:", err)
	  return
	}

	time.Sleep(1000000000 * time.Second)
}
