package udpsender

import (
	"errors"
	"../payload"
	"../session"
)

func SendPayload(currentSession *session.Session, currentPayload *payload.Payload) error {
	udpConn := currentSession.UdpConn
	udpAddr := currentSession.UdpAddr

	if udpConn == nil {
		return errors.New("udpConn is nil")
	}
	if udpAddr == nil {
		return errors.New("udpAddr is nil")
	}

	_, err := udpConn.WriteToUDP(currentPayload.Buffer[:currentPayload.BufferLength], udpAddr)
	if err != nil {
		_, err := udpConn.Write(currentPayload.Buffer[:currentPayload.BufferLength])
		if err != nil {
			return err
		}
	}

	return nil
}

func SendPayloadsForInitialGap(currentSession *session.Session) error {
	sessionId := currentSession.SessionId

	dummyBytes := make([]byte, 1000)

	for i := 0; i < 100; i++ {
		currentPayload := new(payload.Payload)
		currentPayload.Buffer = append([]byte("4GAP"), (*sessionId)[:]...)
		currentPayload.Buffer = append(currentPayload.Buffer, dummyBytes...)
		currentPayload.BufferLength = len(dummyBytes)
		
		SendPayload(currentSession, currentPayload)
	}

	return nil
}
