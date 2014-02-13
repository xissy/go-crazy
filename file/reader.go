package file

import (
	// "fmt"
	"log"
	"os"
	"encoding/binary"
	"errors"
	"time"
	"github.com/nu7hatch/gouuid"
	"github.com/willf/bitset"
	"../udpsender"
	"../payload"
	"../session"
)

const DefaultFileReadingBufferCount int = 5000

func StartToReadFile(sessionId *uuid.UUID, filename string) (*File, error) {
	currentSession, err := session.GetSession(sessionId)
	if err != nil { return nil, err }
	if currentSession == nil { return nil, errors.New("Session is nil: " + sessionId.String()) }

	srcFileHandle, err := os.Open(filename)
	if err != nil { return nil, err }

	fileId, err := uuid.NewV4()
	if err != nil { return nil, err }

	log.Println("started StartToReadFile:", sessionId, filename, fileId)

	payloadChannel := make(chan *payload.Payload, DefaultFileReadingBufferCount)

	file := new(File)
	file.SessionId = sessionId
	file.Session = currentSession
	file.FileId = fileId
	file.SrcFileHandle = srcFileHandle
	file.PayloadDataSize = DefaultPayloadDataSize
	file.PayloadChannel = payloadChannel

	go loopToReadFilePayload(file)
	go loopToPrepareDataPayload(file)
	go loopToSendDataPayload(file)

	return file, nil
}

func loopToReadFilePayload(file *File) error {
	var payloadNo int64
	payloadNo = 0

	for {
		var currentPayload payload.Payload
		currentPayload.UdpAddr = file.Session.UdpAddr
		currentPayload.UdpConn = file.Session.UdpConn
		
		buffer := make([]byte, 1400)
		copy(buffer[:4], []byte("DATA"))
		copy(buffer[4:20], file.FileId[:])
		binary.PutVarint(buffer[20:28], payloadNo)

		bufferLength, err := file.SrcFileHandle.Read(buffer[28:28 + DefaultPayloadDataSize])
		if err != nil {
			file.PayloadChannel <- nil

			return err
		}

		currentPayload.Buffer = buffer
		currentPayload.BufferLength = bufferLength + 28

		file.PayloadChannel <- &currentPayload

		payloadNo++
	}

	defer file.SrcFileHandle.Close()

	return nil
}

func loopToPrepareDataPayload(file *File) error {
	file.SendingPayloadMap = make(map[int64]*payload.Payload, DefaultFileReadingBufferCount)

	for {
		if !file.IsReadingFinished {
			currentPayload := <- file.PayloadChannel

			if currentPayload == nil || currentPayload.Err != nil {
				file.IsReadingFinished = true

				log.Println("ended loopToPrepareDataPayload:", file.FileId)

				return nil
			} else {
				payloadNo, _ := binary.Varint(currentPayload.Buffer[20:28])
				file.SendingPayloadMap[payloadNo] = currentPayload

				// fmt.Println("prepared payloadNo:", payloadNo)
				// fmt.Println("file.IsReadingFinished:", file.IsReadingFinished)
			}
		}
	}

	return nil
}

func loopToSendDataPayload(file *File) error {
	for {
		if file.IsReadingFinished && len(file.SendingPayloadMap) == 0 {
			file.IsSendingFinished = true

			// fmt.Println("len(file.SendingPayloadMap):", len(file.SendingPayloadMap))
			log.Println("ended loopToSendDataPayload:", file.FileId)

			return nil
		}

		sortedPayloadNoList := make(int64array, len(file.SendingPayloadMap))
		i := 0
		for payloadNo, _ := range file.SendingPayloadMap {
			sortedPayloadNoList[i] = payloadNo
			i++
		}
		sortedPayloadNoList.Sort()

		for _, payloadNo := range sortedPayloadNoList {
			gap := file.Session.SendingPayloadGap
			time.Sleep(gap)

			payload := file.SendingPayloadMap[payloadNo]
			// log.Println("payload:", payload)

			// log.Println("file.Session.UdpConn:", file.Session.UdpConn)
			// log.Println("file.Session.UdpAddr:", file.Session.UdpAddr)

			log.Println("sending payloadNo:", payloadNo)
			// fmt.Println("len(file.SendingPayloadMap):", len(file.SendingPayloadMap))

			err := udpsender.SendPayload(file.Session, payload)
			if err != nil {
				log.Println("failed to udp.SendPayload():", err)
			}
		}
	}

	return nil
}

func DeleteFromSendingPayloadMap(file *File, payloadNo int64) error {
	if file.SendingPayloadMap == nil {
		return errors.New("file.SendingPayloadMap is nil.")
	}

	// fmt.Println("len(file.SendingPayloadMap):", len(file.SendingPayloadMap))

	delete(file.SendingPayloadMap, payloadNo)

	return nil
}

func DeleteFromSendingPayloadMapWithBitSet(file *File, startPayloadNo int64, ackBitSet *bitset.BitSet) error {
	if file.SendingPayloadMap == nil {
		return errors.New("file.SendingPayloadMap is nil.")
	}

	var i uint
	for i = 0; i < ackBitSet.Count(); i++ {
		if ackBitSet.Test(i) {
			DeleteFromSendingPayloadMap(file, int64(i))
		}
	}

	return nil
}
