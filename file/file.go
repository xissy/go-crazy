package file

import (
	"os"
	"github.com/nu7hatch/gouuid"
	"github.com/willf/bitset"
	"../session"
)

type File struct {
	SessionId *uuid.UUID
	Session *session.Session

	FileId *uuid.UUID
	SrcFileHandle *os.File
	DestFileHandle *os.File
	SrcFilePath string
	DestFilePath string
	
	FileSize int64
	PayloadDataSize int
	PayloadCountInChunk int
	
	ChunkBufferSize int
	ChunkBuffer []*Chunk
	WaitForChunkBufferSpaceChannel chan bool
	FinishedChunkBitSet *bitset.BitSet

	ReceivingChunkMap map[int]*Chunk
	ReceivedPayloadCount int
	
	IsReadingFinished bool
	IsSendingFinished bool
}
