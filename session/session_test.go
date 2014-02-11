package session

import (
	"testing"
	"fmt"
	"github.com/nu7hatch/gouuid"
)

var _sessionId *uuid.UUID

func TestNewSessionId(t *testing.T) {
	sessionId, err := NewSessionId()
	if err != nil {
		t.Error("Failed to create new sessionId:", err)
	}
	fmt.Println(sessionId)

	_sessionId = sessionId
}

func TestPutSession(t *testing.T) {
	var session Session
	session.SessionId = _sessionId

	err := PutSession(&session)
	if err != nil {
		t.Error("Failed to put a session:", err)
	}
	fmt.Println(session)
}

func TestNewSession(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Error("Failed to create a new session", err)
	}
	fmt.Println(session)
}

func TestGetSession(t *testing.T) {
	session, err := GetSession(_sessionId)
	if err != nil {
		t.Error("Failed to get a session:", err)
	}
	fmt.Println(session)
}
