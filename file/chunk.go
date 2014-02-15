package file

import (
	"log"
	"io"
	"bufio"
	"github.com/nu7hatch/gouuid"
	"github.com/willf/bitset"
	"../payload"
)

const DefaultPayloadCountInChunk int = 1000

type Chunk struct {
	FileId *uuid.UUID
	File *File
	ChunkNo int
	StartPayloadNo int64
	Payloads []*payload.Payload
	PayloadsLength int
	IsLast bool
	ReceivedPayloadBitSet *bitset.BitSet
}

func (file* File) NewChunk(chunkNo int) *Chunk {
	chunk := new(Chunk)

	chunk.FileId = file.FileId
	chunk.File = file
	chunk.ChunkNo = chunkNo
	chunk.StartPayloadNo = int64(chunkNo) * int64(file.PayloadCountInChunk)

	chunk.Payloads = make([]*payload.Payload, file.PayloadCountInChunk)
	
	chunk.ReceivedPayloadBitSet = bitset.New(uint(file.PayloadCountInChunk))

	return chunk
}

func (file *File) ReadChunk(chunkNo int) (*Chunk, error) {
	chunk := file.NewChunk(chunkNo)

	_, err := file.SrcFileHandle.Seek(chunk.StartPayloadNo * int64(file.PayloadDataSize), 0)
	if err != nil { return nil, err }

	currentPayloadNo := chunk.StartPayloadNo
	
	reader := bufio.NewReader(file.SrcFileHandle)
	for i := 0; i < file.PayloadCountInChunk; i++ {
		data := make([]byte, file.PayloadDataSize)

		dataLength, err := reader.Read(data)
		// log.Println("dataLength:", dataLength)
		if err != nil {
			if err == io.EOF {
				chunk.IsLast = true
				break
			}
			
			return nil, err
		}

		currentPayload := NewDataPayload(file, currentPayloadNo, data, dataLength)
		currentPayloadNo++

		chunk.Payloads[i] = currentPayload
		chunk.PayloadsLength++
	}

	return chunk, nil
}

func (chunk *Chunk) WriteToFile(file *File) error {
	_, err := file.DestFileHandle.Seek(chunk.StartPayloadNo * int64(payload.DefaultPayloadDataSize), 0)
	if err != nil { return err }

	writer := bufio.NewWriter(file.DestFileHandle)
	defer writer.Flush()

	for i := 0; i < chunk.PayloadsLength; i++ {
		currentPayload := chunk.Payloads[i]

		_, err := writer.Write(currentPayload.Buffer[28:currentPayload.BufferLength])
		if err != nil { return err }

		currentPayload.Buffer = nil
	}

	log.Println("Written to file:", chunk.ChunkNo)

	chunk = nil

	return nil
}
