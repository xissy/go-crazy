package http

import (
	"testing"
	"fmt"
	"github.com/nu7hatch/gouuid"
)

var _sessionId *uuid.UUID

func TestStartHttpServer(t *testing.T) {
	httpPort := 20080

	err := StartHttpServer(httpPort)
	if err != nil {
		t.Error("Failed to Start HTTP Server")
	}
}

func TestConnect(t *testing.T) {
	baseUrl := "http://127.0.0.1:8080"

	json, err := Connect(baseUrl)
	if err != nil {
		t.Error("Failed to Connect", err)
	}

	fmt.Println(json)

	_sessionId = json.SessionId
}

func TestAuth(t *testing.T) {
	baseUrl := "http://127.0.0.1:8080"
	sessionId := _sessionId

	json, err := Auth(baseUrl, sessionId)
	if err != nil {
		t.Error("Failed to Auth", err)
	}

	fmt.Println(json)
}
