package udp

import (
	"fmt"
	"net"
	"log"
	"encoding/binary"
	"time"
	"errors"
	"github.com/nu7hatch/gouuid"
	"../payload"
	"../session"
)

func loopToReadPayload(payloadChannel chan *payload.Payload,
					udpConn *net.UDPConn) error {
	for {
		buffer := make([]byte, 1400)

		bufferLength, udpAddr, err := udpConn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Failed udpConn.ReadFromUDP():", err)
			continue
		}

		currentPayload := new(payload.Payload)
		currentPayload.UdpAddr = udpAddr
		currentPayload.UdpConn = udpConn
		currentPayload.Buffer = buffer
		currentPayload.BufferLength = bufferLength

		payloadChannel <- currentPayload
	}

	return nil
}

func loopToHandlePayload(payloadChannel chan *payload.Payload) error {
	for {
		currentPayload := <- payloadChannel

		handlePayload(currentPayload)

		currentPayload.Buffer = nil
		currentPayload = nil
	}

	return nil
}

func loopToSendHeartbeat(currentSession *session.Session) {
	for {
		if currentSession.IsDisconnected {
			break
		}

		sessionId := currentSession.SessionId

		buffer := append([]byte("BEAT"), (*sessionId)[:]...)
		currentSession.UdpConn.WriteToUDP(buffer, currentSession.UdpAddr)

		time.Sleep(1 * time.Second)
	}
}

func handlePayload(currentPayload *payload.Payload) error {
	buffer := currentPayload.Buffer
	opCode := string(buffer[:4])

	fmt.Println("BufferLength:", currentPayload.BufferLength)
	fmt.Println("opCode:", opCode)
	// fmt.Println("sessionId or fileId:", sessionId)

	// currentSession, err := session.GetSession(sessionId)
	// if err != nil { return err }
	// currentFile, err := file.GetFile(fileId)
	// if err != nil { return err }
	
	// if currentSession == nil && currentFile == nil {
	// 	return errors.New("currentSession or currentFile is nil: " + sessionId.String())
	// }

	switch opCode {
	case "4GAP":
		sessionId := new(uuid.UUID)
		copy(sessionId[:], buffer[4:20])
		currentSession, err := session.GetSession(sessionId)
		if err != nil { return err }
		if currentSession == nil { return errors.New("currentSession is nil.") }

		if currentSession.UdpConn == nil {
			currentSession.UdpConn = currentPayload.UdpConn
			currentSession.UdpAddr = currentPayload.UdpAddr
			
			go loopToSendHeartbeat(currentSession)
		}

		var t time.Time

		if currentSession.PrevInitialPayloadTime == t {
			currentSession.PrevInitialPayloadTime = time.Now()
		} else {
			now := time.Now()
			gap := now.Sub(currentSession.PrevInitialPayloadTime)
			currentSession.PrevInitialPayloadTime = now
			currentSession.InitialPayloadGapSum += gap
			currentSession.InitialPayloadCount++
			currentSession.InitialPayloadGap = 
				currentSession.InitialPayloadGapSum / 
				time.Duration(currentSession.InitialPayloadCount)
		}

		// fmt.Println(currentSession)
		fmt.Println("currentSession.PrevInitialPayloadTime:", currentSession.PrevInitialPayloadTime)
		fmt.Println("currentSession.InitialPayloadGapSum:", currentSession.InitialPayloadGapSum)
		fmt.Println("currentSession.InitialPayloadCount:", currentSession.InitialPayloadCount)
		fmt.Println("currentSession.InitialPayloadGap:", currentSession.InitialPayloadGap)

	case "BEAT":

	case "DATA":
		fileId := new(uuid.UUID)
		copy(fileId[:], buffer[4:20])
		payloadNo, _ := binary.Varint(buffer[20:28])
		log.Println("received DATA:", payloadNo, fileId)
		
	case "ACK1":

	case "ACK2":
	}

	return nil
}
