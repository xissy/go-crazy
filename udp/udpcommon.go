package udp

import (
	"net"
	"log"
	"os"
	"encoding/binary"
	"time"
	"errors"
	"github.com/nu7hatch/gouuid"
	"github.com/willf/bitset"
	"../payload"
	"../session"
	"../file"
	"../udpsender"
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

		buffer = nil

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
		
		// log.Println("received DATA:",  fileId, payloadNo, currentFile.ReceivedPayloadCount)
		currentFile.ReceivedPayloadCount++

		chunkNo := int(payloadNo / int64(currentFile.PayloadCountInChunk))

		payloadNoInChunk := int(payloadNo % int64(currentFile.PayloadCountInChunk))

		if currentFile.ReceivingChunkMap[chunkNo] == nil {
			chunk := currentFile.NewChunk(chunkNo)

			// calculate PayloadsLength.
			chunkStartPosition := 
				int64(chunkNo) * int64(currentFile.PayloadDataSize) * 
					int64(currentFile.PayloadCountInChunk)
			chunkFullEndPosition := 
				chunkStartPosition + 
					(int64(currentFile.PayloadDataSize) * int64(currentFile.PayloadCountInChunk))
			
			if chunkFullEndPosition <= currentFile.FileSize {
				chunk.PayloadsLength = currentFile.PayloadCountInChunk
			} else {
				chunk.PayloadsLength = 
					int((currentFile.FileSize - chunkStartPosition) / int64(currentFile.PayloadDataSize)) + 1
			}

			chunk.ReceivedPayloadBitSet = bitset.New(uint(chunk.PayloadsLength))

			currentFile.ReceivingChunkMap[chunkNo] = chunk
		}

		chunk := currentFile.ReceivingChunkMap[chunkNo]
		chunk.Payloads[payloadNoInChunk] = currentPayload
		chunk.ReceivedPayloadBitSet.Set(uint(payloadNoInChunk))

		if currentFile.ReceivedPayloadCount % 10 == 0 {
			nak1Payload, _ := chunk.NewNak1Payload()
			if nak1Payload != nil {
				udpsender.SendPayload(nak1Payload)
			}
		}

		if chunk.ReceivedPayloadBitSet.All() {
			nak1Payload, _ := chunk.NewNak1Payload()
			if nak1Payload != nil {
				udpsender.SendPayload(nak1Payload)
			}

			udpsender.SendPayload(nak1Payload)

			if !currentFile.FinishedChunkBitSet.Test(uint(chunkNo)) {
				chunk.WriteToFile(currentFile)

				currentFile.FinishedChunkBitSet.Set(uint(chunkNo))

				if currentFile.FinishedChunkBitSet.All() {
					delete(file.ReceivingFileMap, *fileId)
					// os.Exit(0)
				}
			}
		}
		
	case "NAK1":
		fileId := new(uuid.UUID)
		copy(fileId[:], buffer[4:20])

		currentFile, err := file.SendingFileMap.GetFile(fileId)
		if err != nil { return err }
		if currentFile == nil { return errors.New("currentFile is nil.") }

		chunkNoInt64, _ := binary.Varint(buffer[20:28])		
		chunkNo := int(chunkNoInt64)

		receivedPayloadBitSet := new(bitset.BitSet)
		receivedPayloadBitSet.UnmarshalJSON(buffer[28:bufferLength])

		err = currentFile.DeleteReceivedPayloads(chunkNo, receivedPayloadBitSet)

		if currentFile.FinishedChunkBitSet.All() {
			log.Println("Time:", time.Now().Sub(file.StartToReadFileTime))

			os.Exit(0)
		}

		buffer = nil

	case "ACK1":

	case "ACK2":
	}

	return nil
}
