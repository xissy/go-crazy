package file

import (
	"os"
	"time"
	"github.com/nu7hatch/gouuid"
	"github.com/willf/bitset"
	"../session"
	"../payload"
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

	SendingPayloadMapCapacity int
	SendingPayloadMap map[int64]*payload.Payload
	WaitForSendingPayloadMapSpaceChannel chan bool
	
	ReceivedPayloadBitSet *bitset.BitSet
	AckPayloadData []byte
	LastAck1SentTime time.Time
	
	IsReadingFinished bool
	IsSendingFinished bool
}
