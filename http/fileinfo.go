package http

import (
	"github.com/nu7hatch/gouuid"
)

type FileInfo struct {
	SessionId *uuid.UUID
	FileId *uuid.UUID
	SrcFilePath string
	DestFilePath string
	FileSize int64
	PayloadDataSize int
}
