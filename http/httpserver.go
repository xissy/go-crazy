package http

import (
	"log"
	"fmt"
	"time"
	"net/http"
	"github.com/nu7hatch/gouuid"
	"github.com/ant0ine/go-json-rest"
	"../session"
	"../file"
	"../udpsender"
	"../payload"
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
	currentFile.PayloadChannel = make(chan *payload.Payload)

	file.PutFile(currentFile)

	fmt.Println(file.FileMap)
	fmt.Println(file.FileMap[*currentFile.FileId])

	file.StartToWriteFile(currentFile)

	sessionIdString := req.PathParam("sessionId")
	sessionId, err := uuid.ParseHex(sessionIdString)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Println(sessionId)

	message := Message {
		IsSuccess: true,
		SessionId: sessionId,
	}

	w.WriteJson(&message)	
}

func StartHttpServer(httpPort int) error {
	log.Println("Trying to start HTTP server port:", httpPort)

	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		rest.Route{ "POST", "/sessions", CreateSessionHandler },
		rest.Route{ "GET", "/sessions/:sessionId/auth", AuthSessionHandler },
		rest.Route{ "POST", "/sessions/:sessionId/send", SendFileHandler },
	)

	go http.ListenAndServe(":8080", &handler)
	
	return nil
}
