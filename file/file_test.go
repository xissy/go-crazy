package file

import (
	"testing"
	// "fmt"
	"time"
	"github.com/nu7hatch/gouuid"
)

func TestReadAndWriteFile(t *testing.T) {
	srcFilePath := "../test.mp3"
	destFilePath := "../test.write.mp3"

	sessionId, err := uuid.NewV4()

	srcFile, err := StartToReadFile(sessionId, srcFilePath)
	if err != nil {
	  t.Error("Failed to Start to open a file:", err)
	  return
	}

	destFile := new(File)
	destFile.DestFilePath = destFilePath
	destFile.PayloadChannel = srcFile.PayloadChannel

	err = StartToWriteFile(destFile)

	time.Sleep(1 * time.Second)
}
