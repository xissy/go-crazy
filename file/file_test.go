package file

import (
	"testing"
	// "fmt"
	"time"
	"github.com/nu7hatch/gouuid"
	"../session"
)

func TestReadFile(t *testing.T) {
	srcFilePath := "../test.mp3"

	sessionId, err := uuid.NewV4()
	if err != nil {
	  t.Error("Failed to generate sessionId:", err)
	  return
	}

	currentSession := new(session.Session)
	currentSession.SessionId = sessionId
	session.PutSession(currentSession)

	currentSession.SendingPayloadGap = time.Duration(1000 * time.Nanosecond)

	srcFile, err := StartToReadFile(sessionId, srcFilePath)
	if err != nil {
	  t.Error("Failed to Start to open a file:", err)
	  return
	}

	// fmt.Println(srcFile)

	// go routine for deleting payloads from file.SendingPayloadMap
	go func() {
		var i int64
		i = 0

		time.Sleep(100 * time.Millisecond)

		for {
			DeleteFromSendingPayloadMap(srcFile, i)
			// fmt.Println("DeleteFromSendingPayloadMap:", i)
			i++
			time.Sleep(1000 * time.Nanosecond)

			if i > 10000 {
				i = 0
			}

			if srcFile.IsSendingFinished {
				break
			}
		}
	}()

	time.Sleep(1 * time.Second)
}
