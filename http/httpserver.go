package http

import (
	"log"
	"time"
	"net/http"
	"github.com/nu7hatch/gouuid"
	"github.com/ant0ine/go-json-rest"
	"../session"
	"../file"
	"../udpsender"
)

func CreateSessionHandler(w *rest.ResponseWriter, req *rest.Request) {
	currentSession, err := session.NewSession()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	message := Message {
		IsSuccess: true,
		SessionId: currentSession.SessionId,
	}

	w.WriteJson(&message)
}

func AuthSessionHandler(w *rest.ResponseWriter, req *rest.Request) {
	sessionIdString := req.PathParam("sessionId")
	sessionId, err := uuid.ParseHex(sessionIdString)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	currentSession, err := session.GetSession(sessionId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	isSuccess := false
	if currentSession.UdpAddr != nil {
		udpsender.SendPayloadsForInitialGap(currentSession)
		isSuccess = true

		time.Sleep(100 * time.Millisecond)
	}

	message := Message {
		IsSuccess: isSuccess,
		SessionId: sessionId,
		InitialPayloadGap: currentSession.InitialPayloadGap,
	}

	w.WriteJson(&message)
}

func SendFileHandler(w *rest.ResponseWriter, req *rest.Request) {
	currentFile := new(file.File)
	req.DecodeJsonPayload(currentFile)

	sessionIdString := req.PathParam("sessionId")
	sessionId, err := uuid.ParseHex(sessionIdString)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	currentSession, err := session.GetSession(sessionId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	currentFile.SessionId = currentSession.SessionId
	currentFile.Session = currentSession

	file.ReceivingFileMap.PutFile(currentFile)

	file.StartToWriteFile(currentFile)

	message := Message {
		IsSuccess: true,
		SessionId: sessionId,
	}

	w.WriteJson(&message)	
}

func ReceiveFileHandler(w *rest.ResponseWriter, req *rest.Request) {
	currentFile := new(file.File)
	req.DecodeJsonPayload(currentFile)

	sessionIdString := req.PathParam("sessionId")
	sessionId, err := uuid.ParseHex(sessionIdString)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fileId, err := uuid.NewV4()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	currentFile.FileId = fileId

	currentFile, err = file.StartToReadFile(sessionId, currentFile.FileId, currentFile.SrcFilePath)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Println("currentFile:", currentFile)

	file.SendingFileMap.PutFile(currentFile)	

	message := Message {
		IsSuccess: true,
		SessionId: sessionId,
		FileId: currentFile.FileId,
		FileSize: currentFile.FileSize,
		PayloadDataSize: currentFile.PayloadDataSize,
		PayloadCountInChunk: currentFile.PayloadCountInChunk,
	}

	log.Println("message:", message)

	w.WriteJson(&message)
}

func StartHttpServer(httpPort int) error {
	log.Println("Trying to start HTTP server port:", httpPort)

	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		rest.Route{ "POST", "/sessions", CreateSessionHandler },
		rest.Route{ "GET", "/sessions/:sessionId/auth", AuthSessionHandler },
		rest.Route{ "POST", "/sessions/:sessionId/send", SendFileHandler },
		rest.Route{ "POST", "/sessions/:sessionId/receive", ReceiveFileHandler },
	)

	go http.ListenAndServe(":8080", &handler)
	
	return nil
}
