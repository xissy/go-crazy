package file

import (
	"log"
	"os"
	"errors"
	"time"
	"io"
	"bufio"
	"github.com/nu7hatch/gouuid"
	"../udpsender"
	"../payload"
	"../session"
)

const DefaultSendingPayloadMapCapacity int = 10000

var StartToReadFileTime time.Time

// var incDuration time.Duration

func StartToReadFile(sessionId *uuid.UUID, fileId *uuid.UUID, filename string) (*File, error) {
	currentSession, err := session.GetSession(sessionId)
	if err != nil { return nil, err }
	if currentSession == nil { return nil, errors.New("Session is nil: " + sessionId.String()) }

	srcFileHandle, err := os.Open(filename)
	if err != nil { return nil, err }

	srcFileInfo, err := srcFileHandle.Stat()
	if err != nil { return nil, err }
	srcFileSize := srcFileInfo.Size()

	log.Println("started StartToReadFile:", sessionId, filename, fileId)

	file := new(File)
	file.SessionId = sessionId
	file.Session = currentSession
	file.FileId = fileId
	file.SrcFileHandle = srcFileHandle
	file.FileSize = srcFileSize
	file.PayloadDataSize = payload.DefaultPayloadDataSize
	file.PayloadCountInChunk = DefaultPayloadCountInChunk
	file.SendingPayloadMapCapacity = DefaultSendingPayloadMapCapacity
	file.SendingPayloadMap = make(map[int64]*payload.Payload, file.SendingPayloadMapCapacity)
	file.WaitForSendingPayloadMapSpaceChannel = make(chan bool, file.SendingPayloadMapCapacity)

	SendingFileMap.PutFile(file)

	go file.loopToReadPayload()
	go file.loopToSendPayload()

	return file, nil
}

func (file *File) loopToReadPayload() error {
	srcFileReader := bufio.NewReader(file.SrcFileHandle)
	var payloadNo int64 = 0
	buffer := make([]byte, file.PayloadDataSize)

	for {
		if file.Session.IsDisconnected {
			break
		}

		if len(file.SendingPayloadMap) < file.SendingPayloadMapCapacity {
			bufferLength, err := srcFileReader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					file.IsReadingFinished = true
					break
				} else {
					return err
				}
			}

			// log.Println("srcFileReader.Read(buffer):", bufferLength, payloadNo)

			currentPayload := NewDataPayload(file, payloadNo, buffer, bufferLength)
			file.SendingPayloadMap[payloadNo] = currentPayload

			payloadNo++

		} else {
			if !file.IsReadingFinished {
				<- file.WaitForSendingPayloadMapSpaceChannel
			}
		}
	}

	return nil
}

var prevTime time.Time

func (file *File) loopToSendPayload() error {
	for {
		if file.Session.IsDisconnected {
			break
		}

		sendingPayloadMapLength := len(file.SendingPayloadMap)

		if file.IsReadingFinished && sendingPayloadMapLength == 0 {
			file.IsSendingFinished = true

			log.Println("Ended loopToSendPayload:", file.FileId)
			log.Println("Time:", time.Now().Sub(StartToReadFileTime))
			os.Exit(0)

			return nil
		}

		sortedPayloadNoList := make(int64array, sendingPayloadMapLength)
		i := 0
		for payloadNo, _ := range file.SendingPayloadMap {
			sortedPayloadNoList[i] = payloadNo
			i++
			if i > sendingPayloadMapLength { break }
		}
		sortedPayloadNoList.Sort()

		for i, payloadNo := range sortedPayloadNoList {
			if file.Session.IsDisconnected {
				return nil
			}

			currentPayload := file.SendingPayloadMap[payloadNo]
			if currentPayload == nil {
				continue
			}

			err := udpsender.SendPayload(currentPayload)
			if err != nil {
				log.Println("Failed to udp.SendPayload():", err)
			}

			gap := file.Session.SendingPayloadGap
			// gap -= time.Duration(8 * time.Millisecond)
			// if gap <= time.Duration(0) {
			// 	gap = time.Duration(1 * time.Nanosecond)
			// }
			// log.Println("gap:", gap)
			// if i % 10 == 0 {
			// 	time.Sleep(gap)
			// }
			i = i
			gap = gap
			// time.Sleep(gap - incDuration)
			// incDuration += time.Duration(10 * time.Nanosecond)
			// time.Sleep(gap)

			currentTime := time.Now()
			log.Println("time diff:", currentTime.Sub(prevTime), payloadNo, file.Session.SendingPayloadGap)
			prevTime = currentTime

			// log.Println("send payloadNo:", payloadNo)
		}
	}

	return nil
}
