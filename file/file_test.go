package file

import (
	"testing"
	// "fmt"
	"log"
	"time"
	"github.com/nu7hatch/gouuid"
	"github.com/willf/bitset"
	"../session"
)

func TestReadChunk(t *testing.T) {
	srcFilePath := "../test.mp3"

	sessionId, err := uuid.NewV4()
	if err != nil {
		t.Error("Failed to generate sessionId:", err)
		return
	}

	fileId, err := uuid.NewV4()
	if err != nil {
		t.Error("Failed to generate fileId:", err)
		return
	}

	currentSession := new(session.Session)
	currentSession.SessionId = sessionId
	session.PutSession(currentSession)

	currentSession.SendingPayloadGap = time.Duration(1000 * time.Nanosecond)

	srcFile, err := StartToReadFile(sessionId, fileId, srcFilePath)
	if err != nil {
		t.Error("Failed to Start to open a file:", err)
		return
	}

	srcFile = srcFile

	go func() {
		receivedPayloadBitSet := bitset.New(1000)

		var bitSetIndex uint
		for bitSetIndex = 0; bitSetIndex < receivedPayloadBitSet.Len(); bitSetIndex++ {
			receivedPayloadBitSet = receivedPayloadBitSet.Set(bitSetIndex)
		}

		for i := 0; i < 10; i++ {
			log.Println("DeleteReceivedPayloads called.")
			srcFile.DeleteReceivedPayloads(i, receivedPayloadBitSet)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	time.Sleep(5 * time.Second)
}

// func TestReadFile(t *testing.T) {
// 	srcFilePath := "../test.mp3"

// 	sessionId, err := uuid.NewV4()
// 	if err != nil {
// 		t.Error("Failed to generate sessionId:", err)
// 		return
// 	}

// 	fileId, err := uuid.NewV4()
// 	if err != nil {
// 		t.Error("Failed to generate fileId:", err)
// 		return
// 	}

// 	currentSession := new(session.Session)
// 	currentSession.SessionId = sessionId
// 	session.PutSession(currentSession)

// 	currentSession.SendingPayloadGap = time.Duration(1000 * time.Nanosecond)

// 	srcFile, err := StartToReadFile(sessionId, fileId, srcFilePath)
// 	if err != nil {
// 		t.Error("Failed to Start to open a file:", err)
// 		return
// 	}

// 	// fmt.Println(srcFile)

// 	// go routine for deleting payloads from file.SendingPayloadMap
// 	go func() {
// 		var i int64
// 		i = 0

// 		time.Sleep(100 * time.Millisecond)

// 		for {
// 			DeleteFromSendingPayloadMap(srcFile, i)
// 			// fmt.Println("DeleteFromSendingPayloadMap:", i)
// 			i++
// 			time.Sleep(1000 * time.Nanosecond)

// 			if i > 10000 {
// 				i = 0
// 			}

// 			if srcFile.IsSendingFinished {
// 				break
// 			}
// 		}
// 	}()

// 	time.Sleep(1 * time.Second)
// }
