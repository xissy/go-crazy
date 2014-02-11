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
}
