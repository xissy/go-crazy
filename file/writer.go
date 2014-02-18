package file

import (
	"log"
	"os"
	"time"
	"encoding/binary"
	"github.com/willf/bitset"
	"../udpsender"
)

func StartToWriteFile(file *File) error {
	log.Println("started StartToWriteFile:", file.FileId)

	destFileHandle, err := os.Create(file.DestFilePath)
	if err != nil { return err }

	file.DestFileHandle = destFileHandle

	receivedPayloadBitSetCapacity := uint(file.FileSize / int64(file.PayloadDataSize))
	if file.FileSize % int64(file.PayloadDataSize) > 0 {
		receivedPayloadBitSetCapacity++
	}
	file.ReceivedPayloadBitSet = bitset.New(receivedPayloadBitSetCapacity)

	file.AckPayloadData = make([]byte, 0, file.PayloadDataSize)

	return nil
}

func (file *File)BuildAndSendAckPayload(receivedPayloadNo int64) error {
	log.Println("BuildAndSendAckPayload(receivedPayloadNo):", receivedPayloadNo)

	ackPayloadLength := len(file.AckPayloadData)
	log.Println("ackPayloadLength:", ackPayloadLength)

	var lastPayloadNo int64
	if ackPayloadLength >= 16 {
		lastPayloadNo, _ = binary.Varint(file.AckPayloadData[ackPayloadLength-8:])
	} else {
		lastPayloadNo = 0
	}

	if ackPayloadLength >= 16 && receivedPayloadNo == lastPayloadNo + 1 {
		binary.PutVarint(file.AckPayloadData[ackPayloadLength-8:], receivedPayloadNo)
	} else {
		receivedPayloadNoPairBytes := make([]byte, 16)
		binary.PutVarint(receivedPayloadNoPairBytes[:8], receivedPayloadNo)
		binary.PutVarint(receivedPayloadNoPairBytes[8:], receivedPayloadNo)
		file.AckPayloadData = append(file.AckPayloadData, receivedPayloadNoPairBytes...)
	}

	if len(file.AckPayloadData) + 16 > file.PayloadDataSize {
		ack1Payload := file.NewAck1Payload(file.AckPayloadData)
		udpsender.SendPayload(ack1Payload)
		file.AckPayloadData = file.AckPayloadData[:0]

		// log.Println("ACK1 sent:", ack1Payload)

		file.LastAck1SentTime = time.Now()
	} else {
		currentTime := time.Now()
		if currentTime.Sub(file.LastAck1SentTime) > time.Duration(10 * time.Millisecond) {
			ack1Payload := file.NewAck1Payload(file.AckPayloadData)
			udpsender.SendPayload(ack1Payload)
			// log.Println("err := udpsender.SendPayload(ack1Payload):", err)
			file.AckPayloadData = file.AckPayloadData[:0]

			// log.Println("ACK1 sent:", ack1Payload)

			file.LastAck1SentTime = time.Now()
		}
	}

	return nil
}
