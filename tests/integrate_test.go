package crazy

import (
	"testing"
	"fmt"
	"os"
	"time"
	"github.com/nu7hatch/gouuid"
	"../http"
	// "../udp"
	"../server"
	"../client"
	"../file"
)

func TestInitialPayloadGap(t *testing.T) {
	host := "127.0.0.1"
	// host := "dev.pricepoller.com"
	httpPort := 8080
	tcpPort := 20068
	udpPort := 20069
	baseUrl := fmt.Sprintf("http://%s:%d", host, httpPort)

	err := server.StartServer(httpPort, tcpPort, udpPort)
	if err != nil {
	  t.Error("Failed to Start the Server:", err)
	  return
	}

	currentSession, err := client.ConnectToServer(host, httpPort, tcpPort, udpPort)
	if err != nil {
	  t.Error("Failed to Start Server:", err)
	  return
	}

	sessionId := currentSession.SessionId

	fileId, err := uuid.NewV4()

	srcFileInfo, err := os.Stat("../test.mp3")
	if err != nil {
		t.Error("Failed to get srcFileInfo:", err)
	}
	srcFileSize := srcFileInfo.Size()
	fmt.Println("srcFileSize:", srcFileSize)

	destFileInfo := new(http.FileInfo)
	destFileInfo.SessionId = sessionId
	destFileInfo.FileId = fileId
	destFileInfo.DestFilePath = "./test.write.mp3"
	destFileInfo.FileSize = srcFileSize
	destFileInfo.PayloadDataSize = file.DefaultPayloadDataSize

	sendFileJson, err := http.SendFile(baseUrl, sessionId, destFileInfo)
	if err != nil {
		t.Error("Failed to Send File:", err)
		return
	}
	fmt.Println("sendFileJson:", sendFileJson)

	srcFile, err := file.StartToReadFile(sessionId, "../test.mp3")
	if err != nil {
		t.Error("Failed to Start to read the source file:", err)
		return
	}
	fmt.Println("srcFile:", srcFile)


	time.Sleep(1000000 * time.Second)
}
