package udp

import (
	// "fmt"
	"net"
	"log"
	// "encoding/binary"
	"time"
	"errors"
	"github.com/nu7hatch/gouuid"
	"../session"
)

func loopToReadPayload(payloadChannel chan *Payload,
					udpConn *net.UDPConn) error {
	for {
		buffer := make([]byte, 1400)

		bufferLength, udpAddr, err := udpConn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Failed udpConn.ReadFromUDP():", err)
			continue
		}

		payload := new(Payload)
		payload.UdpAddr = udpAddr
		payload.UdpConn = udpConn
		payload.Buffer = buffer
		payload.BufferLength = bufferLength

		payloadChannel <- payload
	}

	return nil
}

func loopToHandlePayload(payloadChannel chan *Payload) error {
	for {
		payload := <- payloadChannel

		handlePayload(payload)

		payload.Buffer = nil
		payload = nil
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

func handlePayload(payload *Payload) error {
	buffer := payload.Buffer
	opCode := string(buffer[:4])
	sessionId := new(uuid.UUID)
	copy(sessionId[:], buffer[4:20])

	// fmt.Println("BufferLength:", payload.BufferLength)
	// fmt.Println("opCode:", opCode)
	// fmt.Println("sessionId:", sessionId)

	currentSession, err := session.GetSession(sessionId)
	if err != nil { return err }
	if currentSession == nil { return errors.New("Session is nil: " + sessionId.String()) }

	switch opCode {
	case "4GAP":
		if currentSession.UdpConn == nil {
			currentSession.UdpConn = payload.UdpConn
			currentSession.UdpAddr = payload.UdpAddr
			
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
		// fmt.Println("currentSession.PrevInitialPayloadTime:", currentSession.PrevInitialPayloadTime)
		// fmt.Println("currentSession.InitialPayloadGapSum:", currentSession.InitialPayloadGapSum)
		// fmt.Println("currentSession.InitialPayloadCount:", currentSession.InitialPayloadCount)
		// fmt.Println("currentSession.InitialPayloadGap:", currentSession.InitialPayloadGap)

	case "BEAT":

	case "DATA":
		
	}

	return nil
}

func SendPayload(currentSession *session.Session, payload *Payload) error {
	udpConn := currentSession.UdpConn
	udpAddr := currentSession.UdpAddr

	_, err := udpConn.WriteToUDP(payload.Buffer[:payload.BufferLength], udpAddr)
	if err != nil {
		_, err := udpConn.Write(payload.Buffer[:payload.BufferLength])
		if err != nil {
			return err
		}
	}

	return nil
}

func SendPayloadsForInitialGap(currentSession *session.Session) error {
	sessionId := currentSession.SessionId

	dummyBytes := make([]byte, 1000)

	for i := 0; i < 100; i++ {
		payload := new(Payload)
		payload.Buffer = append([]byte("4GAP"), (*sessionId)[:]...)
		payload.Buffer = append(payload.Buffer, dummyBytes...)
		payload.BufferLength = len(dummyBytes)
		
		SendPayload(currentSession, payload)
	}

	return nil
}
