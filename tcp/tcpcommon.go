package tcp

import (
	"log"
	"net"
	"errors"
	"github.com/nu7hatch/gouuid"
	"../session"
)

func loopToReadStream(streamMap map[*net.TCPConn]*Stream,
					streamChannel chan *Stream,
					tcpConn *net.TCPConn) error {
	for {
		buffer := make([]byte, 8000)
		bufferLength, err := tcpConn.Read(buffer)
		if err != nil {
			currentStream := streamMap[tcpConn]
			sessionId := currentStream.SessionId

			currentSession, err := session.GetSession(sessionId)
			if currentSession != nil {
				// set that this session is disconnected to stop heartbeat loop.
				currentSession.IsDisconnected = true

				session.DeleteSession(sessionId)

				log.Println("Disconnected.")

				// TODO: delete the files in FileMap.
			}

			return err
		}

		stream := new(Stream)
		stream.TcpConn = tcpConn
		stream.Buffer = buffer
		stream.BufferLength = bufferLength

		streamChannel <- stream
	}

	defer tcpConn.Close()
	defer cleanStream(streamMap, tcpConn)

	return nil
}

func cleanStream(streamMap map[*net.TCPConn]*Stream, tcpConn *net.TCPConn) {
	delete(streamMap, tcpConn)
}

func loopToHandleStream(streamMap map[*net.TCPConn]*Stream,
					streamChannel chan *Stream) error {
	for {
		stream := <- streamChannel

		// if this stream is the first of the TcpConn then
		_, isInMap := streamMap[stream.TcpConn]
		if !isInMap {
			// put the stream at streamMap and 
			streamMap[stream.TcpConn] = stream
			// **write same stream back.**
			stream.TcpConn.Write(stream.Buffer[:stream.BufferLength])
		}

		buffer := streamMap[stream.TcpConn].Buffer
		buffer = append(buffer, stream.Buffer[:stream.BufferLength]...)
		streamMap[stream.TcpConn].BufferLength += stream.BufferLength
		
		sessionId, err := getSessionIdFromStream(buffer)
		if err == nil {
			currentStream := streamMap[stream.TcpConn]
			currentStream.SessionId = sessionId

			currentSession, err := session.GetSession(sessionId)
			if currentSession != nil && err == nil {
				currentSession.TcpConn = stream.TcpConn
			}
		}

		stream.Buffer = nil
		stream = nil
	}

	return nil
}

func getSessionIdFromStream(stream []byte) (*uuid.UUID, error) {
	if string(stream[0:4]) != "AUTH" {
		return nil, errors.New("invalid stream. it doesn't start with AUTH")
	}
	if len(stream) < 20 {
		return nil, errors.New("stream length is too short. should be 20 or more")
	}

	sessionId, err := uuid.Parse(stream[4:20])
	if err != nil { return nil, err }

	return sessionId, nil
}
