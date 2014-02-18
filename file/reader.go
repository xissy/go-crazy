package file

import (
	"log"
	"os"
	"errors"
	"time"
	"github.com/nu7hatch/gouuid"
	"github.com/willf/bitset"
	"../udpsender"
	"../payload"
	"../session"
)

const DefaultChunkBufferSize int = 10

var StartToReadFileTime time.Time

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
	file.PayloadDataSize = payload.DefaultPayloadDataSize
	file.PayloadCountInChunk = DefaultPayloadCountInChunk
	file.ChunkBufferSize = DefaultChunkBufferSize
	file.WaitForChunkBufferSpaceChannel = make(chan bool)
	chunkCount := uint(srcFileSize  / int64(file.PayloadDataSize) / int64(file.PayloadCountInChunk)) + uint(1)
	file.FinishedChunkBitSet = bitset.New(chunkCount)

	SendingFileMap.PutFile(file)

	go file.loopToReadChunk()
	go file.loopToSendDataPayload()

	return file, nil
}

func (file *File) loopToReadChunk() error {
	file.ChunkBuffer = make([]*Chunk, 0, file.ChunkBufferSize)
	currentChunkNo := 0

	for {
		if file.Session.IsDisconnected {
			break
		}

		if len(file.ChunkBuffer) < file.ChunkBufferSize {
			currentChunk, err := file.ReadChunk(currentChunkNo)
			if err != nil { return err }

			log.Println("ReadChunk:", currentChunkNo)

			file.ChunkBuffer = append(file.ChunkBuffer, currentChunk)

			if currentChunk.IsLast {
				file.IsReadingFinished = true
				break
			}

			currentChunkNo++
		} else {
			if !file.IsReadingFinished {
				<- file.WaitForChunkBufferSpaceChannel
			}
		}
	}

	return nil
}

var prevTime time.Time
var i int

func (file *File) loopToSendDataPayload() error {
	for {
		if file.Session.IsDisconnected { break }

		if file.FinishedChunkBitSet.All() { break }

		sendingPayloadsCapacity := len(file.ChunkBuffer) * file.PayloadCountInChunk
		sendingPayloads := make([]*payload.Payload, 0, sendingPayloadsCapacity)
		for _, currentChunk := range file.ChunkBuffer {
			for payloadPos, currentPayload := range currentChunk.Payloads[:currentChunk.PayloadsLength] {
				if !currentChunk.ReceivedPayloadBitSet.Test(uint(payloadPos)) {
					sendingPayloads = append(sendingPayloads, currentPayload)
				}
			}
		}

		time.Sleep(1 * time.Nanosecond)

		for _, currentPayload := range sendingPayloads {
			if file.Session.IsDisconnected { return nil }

			err := udpsender.SendPayload(currentPayload)
			if err != nil {
				log.Println("failed to udp.SendPayload():", err)
			}

			gap := file.Session.SendingPayloadGap
			// log.Println("gap:", gap)
			// gap = gap / 2
			gap -= time.Duration(16 * time.Microsecond)
			// gap = time.Duration(120000 * time.Nanosecond)
			if gap < 0 {
				gap = time.Duration(1 * time.Nanosecond)
			}
			// log.Println("gap:", gap)
			if i % 10 == 0 {
				time.Sleep(gap)
			}
			// time.Sleep(gap)
			i++

			currentTime := time.Now()
			log.Println("time diff:", currentTime.Sub(prevTime), file.Session.SendingPayloadGap)
			prevTime = currentTime
		}
	}

	return nil
}

func (file *File) DeleteReceivedPayloads(chunkNo int, receivedPayloadBitSet *bitset.BitSet) error {
	for index, currentChunk := range file.ChunkBuffer {
		if currentChunk.ChunkNo == chunkNo {
			if receivedPayloadBitSet.All() {
				file.ChunkBuffer = append(file.ChunkBuffer[:index], file.ChunkBuffer[index+1:]...)
				file.FinishedChunkBitSet.Set(uint(chunkNo))

				if !file.IsReadingFinished {
					file.WaitForChunkBufferSpaceChannel <- true
				}

			} else {
				currentChunk.ReceivedPayloadBitSet = receivedPayloadBitSet
			}
			break
		}
	}

	return nil
}
