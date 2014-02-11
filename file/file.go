package file

import (
	"os"
	"github.com/nu7hatch/gouuid"
	"../udp"
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
	PayloadChannel chan *udp.Payload

	// Chunks []*Chunk
}
