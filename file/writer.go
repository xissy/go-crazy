package file

import (
	"fmt"
	"log"
	"os"
	"encoding/binary"
	"github.com/willf/bitset"
)

func StartToWriteFile(file *File) error {
	log.Println("started StartToWriteFile:", file.FileId)

	destFileHandle, err := os.Create(file.DestFilePath)
	if err != nil {
		return err
	}

	file.DestFileHandle = destFileHandle
	
	bitSetCount := uint(file.FileSize / int64(file.PayloadDataSize))
	file.ReceivedPayloadBitSet = bitset.New(bitSetCount)
	fmt.Println("bitSetCount:", bitSetCount)

	go loopToWritePayloadToFile(file)

	return nil
}

func loopToWritePayloadToFile(file *File) error {
	for {
		payload := <- file.PayloadChannel
		// TODO: exit the loop when file writing is done.
		if payload == nil || payload.Err != nil {
			break
		}

		payloadNo, _ := binary.Varint(payload.Buffer[20:28])

		file.DestFileHandle.Seek(payloadNo * int64(DefaultPayloadDataSize), 0)
		file.DestFileHandle.Write(payload.Buffer[28:28 + payload.BufferLength])

		fmt.Println("writing to file, payloadNo:", payloadNo, ", payload.BufferLength:", payload.BufferLength)
	}

	defer file.DestFileHandle.Close()

	log.Println("ended loopToWritePayloadToFile:", file.FileId)

	return nil
}
