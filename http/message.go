package http

import (
	"time"
	"github.com/nu7hatch/gouuid"
)

type Message struct {
	IsSuccess bool
	SessionId *uuid.UUID
	
	InitialPayloadGap time.Duration

	FileId *uuid.UUID
	SrcFilePath string
	DestFilePath string
	FileSize int64
	PayloadDataSize int
	PayloadCountInChunk int
}
