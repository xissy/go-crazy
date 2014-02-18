package file

import (
	"encoding/binary"
	"../payload"
)

func NewDataPayload(file *File, payloadNo int64, data []byte, dataLength int) *payload.Payload {
	currentPayload := new(payload.Payload)

	currentPayload.UdpAddr = file.Session.UdpAddr
	currentPayload.UdpConn = file.Session.UdpConn

	bufferLength := 28 + dataLength
	buffer := make([]byte, bufferLength)
	copy(buffer[0:4], "DATA")
	copy(buffer[4:20], file.FileId[:])
	binary.PutVarint(buffer[20:28], payloadNo)
	copy(buffer[28:], data[:dataLength])

	currentPayload.Buffer = buffer
	currentPayload.BufferLength = bufferLength

	return currentPayload
}

func (chunk *Chunk) NewNak1Payload() (*payload.Payload, error) {
	currentPayload := new(payload.Payload)

	currentPayload.UdpAddr = chunk.File.Session.UdpAddr
	currentPayload.UdpConn = chunk.File.Session.UdpConn

	buffer := make([]byte, 0, 1400)
	buffer = append(buffer, []byte("NAK1")...)
	buffer = append(buffer, chunk.File.FileId[:]...)

	chunkNoBytes := make([]byte, 8)
	binary.PutVarint(chunkNoBytes, int64(chunk.ChunkNo))
	buffer = append(buffer, chunkNoBytes...)

	receivedPayloadBitSetBytes, err := chunk.ReceivedPayloadBitSet.MarshalJSON()
	if err != nil { return nil, err }
	buffer = append(buffer, receivedPayloadBitSetBytes...)

	currentPayload.Buffer = buffer
	currentPayload.BufferLength = len(buffer)

	return currentPayload, nil
}

func (file *File) NewAck1Payload(ackData []byte) *payload.Payload {
	currentPayload := new(payload.Payload)

	currentPayload.UdpAddr = file.Session.UdpAddr
	currentPayload.UdpConn = file.Session.UdpConn

	buffer := make([]byte, 0, 1400)
	buffer = append(buffer, []byte("ACK1")...)
	buffer = append(buffer, file.FileId[:]...)
	buffer = append(buffer, ackData...)

	currentPayload.Buffer = buffer
	currentPayload.BufferLength = len(buffer)

	return currentPayload
}
