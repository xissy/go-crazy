package file

import (
	"os"
	"github.com/nu7hatch/gouuid"
	"github.com/willf/bitset"
	"../payload"
	"../session"
)

const DefaultPayloadDataSize int = 1000

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
	PayloadChannel chan *payload.Payload
	SendingPayloadMap map[int64]*payload.Payload
	IsReadingFinished bool
	IsSendingFinished bool
	ReceivedPayloadBitSet *bitset.BitSet

	// Chunks []*Chunk
}
