package session

import (
	"errors"
	"github.com/nu7hatch/gouuid"
)

var SessionMap map[uuid.UUID]*Session

func checkAndNewSessionMap() {
	if SessionMap == nil {
		SessionMap = make(map[uuid.UUID]*Session)
	}
}

func NewSessionId() (*uuid.UUID, error) {
	sessionId, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return sessionId, nil
}

func PutSession(session *Session) error {
	checkAndNewSessionMap()

	if session.SessionId == nil {
		return errors.New("invalid sessionId")
	}

	SessionMap[*session.SessionId] = session

	return nil
}

func NewSession() (*Session, error) {
	sessionId, err := NewSessionId()
	if err != nil { return nil, err }

	session := new(Session)
	session.SessionId = sessionId
	err = PutSession(session)
	if err != nil { return nil, err }

	return session, nil
}

func GetSession(sessionId *uuid.UUID) (*Session, error) {
	checkAndNewSessionMap()

	return SessionMap[*sessionId], nil
}

func DeleteSession(sessionId *uuid.UUID) {
	checkAndNewSessionMap()	

	delete(SessionMap, *sessionId)
}
