package crazy

import (
	"log"
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
	"../payload"
)

func TestSendingFile(t *testing.T) {
	host := "127.0.0.1"
	// host := "dev.pricepoller.com"
	// host := "128.199.245.227"
	// host := "track.recom.io"
	httpPort := 8080
	tcpPort := 20068
	udpPort := 20069
	baseUrl := fmt.Sprintf("http://%s:%d", host, httpPort)

	file.SendingFileMap = make(file.FileMap)
	file.ReceivingFileMap = make(file.FileMap)

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

	log.Println("SendingPayloadGap:", currentSession.SendingPayloadGap)
	log.Println("ReceivingPayloadGap:", currentSession.ReceivingPayloadGap)

	sessionId := currentSession.SessionId

	fileId, err := uuid.NewV4()

	srcFilePath := "../test.mp3"
	// srcFilePath := "../jdk.tar"
	srcFileInfo, err := os.Stat(srcFilePath)
	if err != nil {
		t.Error("Failed to get srcFileInfo:", err)
	}
	srcFileSize := srcFileInfo.Size()

	destFileInfo := new(http.FileInfo)
	destFileInfo.SessionId = sessionId
	destFileInfo.FileId = fileId
	destFileInfo.DestFilePath = "./test.write.mp3"
	destFileInfo.FileSize = srcFileSize
	destFileInfo.PayloadDataSize = payload.DefaultPayloadDataSize
	destFileInfo.PayloadCountInChunk = file.DefaultPayloadCountInChunk

	sendFileJson, err := http.SendFile(baseUrl, sessionId, destFileInfo)
	if err != nil {
		t.Error("Failed to Send File:", err)
		return
	}
	sendFileJson = sendFileJson

	file.StartToReadFileTime = time.Now()

	srcFile, err := file.StartToReadFile(sessionId, fileId, srcFilePath)
	if err != nil {
		t.Error("Failed to Start to read the source file:", err)
		return
	}
	srcFile = srcFile

	time.Sleep(1000000 * time.Second)
}

// func TestReceivingFile(t *testing.T) {
// 	host := "127.0.0.1"
// 	// host := "dev.pricepoller.com"
// 	// host := "128.199.245.227"
// 	// host := "track.recom.io"
// 	httpPort := 8080
// 	tcpPort := 20068
// 	udpPort := 20069
// 	baseUrl := fmt.Sprintf("http://%s:%d", host, httpPort)

// 	file.SendingFileMap = make(file.FileMap)
// 	file.ReceivingFileMap = make(file.FileMap)

// 	err := server.StartServer(httpPort, tcpPort, udpPort)
// 	if err != nil {
// 	  t.Error("Failed to Start the Server:", err)
// 	  return
// 	}

// 	currentSession, err := client.ConnectToServer(host, httpPort, tcpPort, udpPort)
// 	if err != nil {
// 	  t.Error("Failed to Start Server:", err)
// 	  return
// 	}

// 	log.Println("SendingPayloadGap:", currentSession.SendingPayloadGap)
// 	log.Println("ReceivingPayloadGap:", currentSession.ReceivingPayloadGap)

// 	sessionId := currentSession.SessionId

// 	srcFileInfo := new(http.FileInfo)
// 	srcFileInfo.SessionId = sessionId
// 	srcFileInfo.SrcFilePath = "../test.mp3"
// 	// srcFileInfo.SrcFilePath = "../jdk.tar"

// 	receivingFileJson, err := http.ReceiveFile(baseUrl, sessionId, srcFileInfo)
// 	if err != nil {
// 		t.Error("Failed to Receiving File:", err)
// 		return
// 	}

// 	log.Println("receivingFileJson:", receivingFileJson)

// 	destFile := new(file.File)
// 	destFile.SessionId = sessionId
// 	destFile.Session = currentSession
// 	destFile.FileId = receivingFileJson.FileId
// 	destFile.DestFilePath = "./test.write.mp3"
// 	destFile.FileSize = receivingFileJson.FileSize
// 	destFile.PayloadDataSize = receivingFileJson.PayloadDataSize
// 	destFile.PayloadCountInChunk = receivingFileJson.PayloadCountInChunk

// 	file.ReceivingFileMap.PutFile(destFile)

// 	file.StartToReadFileTime = time.Now()

// 	err = file.StartToWriteFile(destFile)
// 	if err != nil {
// 		t.Error("Failed to Start to write the file:", err)
// 		return
// 	}

// 	time.Sleep(1000000 * time.Second)
// }
