package udpsender

import (
	"errors"
	"../payload"
	"../session"
)

func SendPayload(currentPayload *payload.Payload) error {
	udpAddr := currentPayload.UdpAddr
	udpConn := currentPayload.UdpConn

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

	dummyBytes := make([]byte, payload.DefaultPayloadDataSize)

	for i := 0; i < 100; i++ {
		currentPayload := new(payload.Payload)

		currentPayload.UdpAddr = currentSession.UdpAddr
		currentPayload.UdpConn = currentSession.UdpConn

		currentPayload.Buffer = append([]byte("4GAP"), (*sessionId)[:]...)
		currentPayload.Buffer = append(currentPayload.Buffer, dummyBytes...)
		currentPayload.BufferLength = len(dummyBytes)
		
		SendPayload(currentPayload)
	}

	return nil
}
