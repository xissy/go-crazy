package main

import (
	"fmt"
	"time"
	"./server"
)

func main() {
	httpPort := 8080
	tcpPort := 20068
	udpPort := 20069

	err := server.StartServer(httpPort, tcpPort, udpPort)
	if err != nil {
	  fmt.Println("Failed to Start the Server:", err)
	  return
	}

	time.Sleep(1000000000 * time.Second)
}
