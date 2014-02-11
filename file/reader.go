package file

import (
	"os"
	"encoding/binary"
	"errors"
	"github.com/nu7hatch/gouuid"
	"../udp"
	"../session"
)

func StartToReadFile(sessionId *uuid.UUID, filename string) (*File, error) {
	currentSession, err := session.GetSession(sessionId)
	if err != nil { return nil, err }
	if currentSession == nil { return nil, errors.New("Session is nil: " + sessionId.String()) }

	srcFileHandle, err := os.Open(filename)
	if err != nil { return nil, err }

	fileId, err := uuid.NewV4()
	if err != nil { return nil, err }

	payloadChannel := make(chan *udp.Payload)

	file := new(File)
	file.SessionId = sessionId
	file.Session = currentSession
	file.FileId = fileId
	file.SrcFileHandle = srcFileHandle
	file.PayloadDataSize = DefaultPayloadDataSize
	file.PayloadChannel = payloadChannel

	go loopToReadFilePayloads(file)

	return file, nil
}

func loopToReadFilePayloads(file *File) error {
	var payloadNo int64
	payloadNo = 0

	for {
		var payload udp.Payload
		
		buffer := make([]byte, 1400)
		copy(buffer[:4], []byte("DATA"))
		copy(buffer[4:20], file.FileId[:])
		binary.PutVarint(buffer[20:28], payloadNo)

		bufferLength, err := file.SrcFileHandle.Read(buffer[28:28 + DefaultPayloadDataSize])
		if err != nil {
			return err
		}

		payload.Buffer = buffer
		payload.BufferLength = bufferLength

		file.PayloadChannel <- &payload

		payloadNo++
	}

	defer file.SrcFileHandle.Close()

	return nil
}
