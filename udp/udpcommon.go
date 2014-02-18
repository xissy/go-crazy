package udp

import (
	"net"
	"log"
	"encoding/binary"
	"time"
	"errors"
	"github.com/nu7hatch/gouuid"
	"../payload"
	"../session"
	"../file"
)

var i int

func loopToReadPayload(payloadChannel chan *payload.Payload,
					udpConn *net.UDPConn) error {
	i++

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
	bufferLength := currentPayload.BufferLength
	opCode := string(buffer[:4])

	defer func() {
		buffer = nil
		currentPayload = nil
	}()

	// log.Println("BufferLength:", currentPayload.BufferLength)
	// log.Println("opCode:", opCode)

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
			if currentSession.InitialPayloadCount > 100 {
				currentSession.InitialPayloadGap = 
					(currentSession.InitialPayloadGapSum - time.Duration(100 * time.Millisecond))/ 
					time.Duration(currentSession.InitialPayloadCount)
			} else {
				currentSession.InitialPayloadGap = 
					currentSession.InitialPayloadGapSum / 
					time.Duration(currentSession.InitialPayloadCount)
			}
		}

		// log.Println("currentSession.PrevInitialPayloadTime:", currentSession.PrevInitialPayloadTime)
		// log.Println("currentSession.InitialPayloadGapSum:", currentSession.InitialPayloadGapSum)
		// log.Println("currentSession.InitialPayloadCount:", currentSession.InitialPayloadCount)
		// log.Println("currentSession.InitialPayloadGap:", currentSession.InitialPayloadGap)

	case "BEAT":

	case "DATA":
		fileId := new(uuid.UUID)
		copy(fileId[:], buffer[4:20])

		currentFile, err := file.ReceivingFileMap.GetFile(fileId)
		if err != nil { return err }
		if currentFile == nil { return errors.New("currentFile is nil.") }

		payloadNo, _ := binary.Varint(buffer[20:28])
		payloadNo = payloadNo
		
		log.Println("received DATA:",  fileId, payloadNo, i)
		i++

		currentFile.BuildAndSendAckPayload(payloadNo)

		if !currentFile.ReceivedPayloadBitSet.Test(uint(payloadNo)) {
			var payloadOffset int64 = payloadNo * int64(currentFile.PayloadDataSize)
			currentFile.DestFileHandle.WriteAt(buffer[28:bufferLength], payloadOffset)

			currentFile.ReceivedPayloadBitSet.Set(uint(payloadNo))
		}
		
	case "NAK1":

	case "ACK1":
		fileId := new(uuid.UUID)
		copy(fileId[:], buffer[4:20])

		currentFile, err := file.SendingFileMap.GetFile(fileId)
		if err != nil { return err }
		if currentFile == nil { return errors.New("currentFile is nil.") }

		ackData := buffer[20:bufferLength]
		// log.Println("ackData:", ackData)
		for pos := 0; pos < len(ackData); pos += 16 {
			startPayloadNo, _ := binary.Varint(ackData[pos:pos+8])
			endPayloadNo, _ := binary.Varint(ackData[pos+8:])
			// log.Println("startPayloadNo, endPayloadNo:", startPayloadNo, endPayloadNo)

			for payloadNo := startPayloadNo; payloadNo <= endPayloadNo; payloadNo++ {
				// log.Println("payloadNo:", payloadNo)
				delete(currentFile.SendingPayloadMap, payloadNo)

				if !currentFile.IsReadingFinished {
					currentFile.WaitForSendingPayloadMapSpaceChannel <- true
				}
			}
		}

		// log.Println("After ACK1:", len(currentFile.SendingPayloadMap))
		log.Println("After ACK1:")

	case "ACK2":
	}

	return nil
}
