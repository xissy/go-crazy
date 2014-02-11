package file

import (
	"github.com/nu7hatch/gouuid"
	"../udp"
)

type Chunk struct {
	SessionId *uuid.UUID
	Payloads []*udp.Payload
	StartPosition uint64
}
