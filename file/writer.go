package file

import (
	"os"
	"fmt"
	"encoding/binary"
)

func StartToWriteFile(file *File) error {
	destFileHandle, err := os.Create(file.DestFilePath)
	if err != nil {
		return err
	}

	file.DestFileHandle = destFileHandle

	go loopToWritePayloadToFile(file)

	return nil
}

func loopToWritePayloadToFile(file *File) error {
	for {
		payload := <- file.PayloadChannel
		// TODO: exit the loop when file writing is done.
		if payload.Err != nil {
			break
		}

		payloadNo, _ := binary.Varint(payload.Buffer[20:28])

		file.DestFileHandle.Seek(payloadNo * int64(DefaultPayloadDataSize), 0)
		file.DestFileHandle.Write(payload.Buffer[28:28 + payload.BufferLength])

		fmt.Println("writing to file, payloadNo:", payloadNo, ", payload.BufferLength:", payload.BufferLength)
	}

	defer file.DestFileHandle.Close()

	return nil
}
