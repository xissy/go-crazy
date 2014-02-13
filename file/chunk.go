package file

import (
	"github.com/nu7hatch/gouuid"
	"../payload"
)

type Chunk struct {
	SessionId *uuid.UUID
	Payloads []*payload.Payload
	StartPosition uint64
}
