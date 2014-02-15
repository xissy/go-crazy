package file

import (
	"log"
	"os"
	"github.com/willf/bitset"
)

func StartToWriteFile(file *File) error {
	log.Println("started StartToWriteFile:", file.FileId)

	destFileHandle, err := os.Create(file.DestFilePath)
	if err != nil { return err }

	file.DestFileHandle = destFileHandle
	file.ReceivingChunkMap = make(map[int]*Chunk)
	chunkCount := uint(file.FileSize / int64(file.PayloadDataSize) / int64(file.PayloadCountInChunk) + 1)
	file.FinishedChunkBitSet = bitset.New(chunkCount)

	return nil
}
