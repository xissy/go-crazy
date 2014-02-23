package http

import (
	"fmt"
	"log"
	"bytes"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/nu7hatch/gouuid"
)

func unmarshalResponse(res *http.Response) (*Message, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil { return nil, err }
	res.Body.Close()

	var message Message
	json.Unmarshal(body, &message)

	return &message, nil
}

func Connect(baseUrl string) (*Message, error) {
	log.Println("Trying to connect to HTTP server")

	finalUrl := baseUrl + "/sessions"

	res, err := http.Post(finalUrl, "", nil)
	if err != nil { return nil, err }
	
	message, err := unmarshalResponse(res)
	if err != nil { return nil, err }

	return message, nil
}

func Auth(baseUrl string, sessionId *uuid.UUID) (*Message, error) {
	finalUrl := fmt.Sprintf("%s/sessions/%s/auth", baseUrl, sessionId.String())

	res, err := http.Get(finalUrl)
	if err != nil { return nil, err }
	
	message, err := unmarshalResponse(res)
	if err != nil { return nil, err }

	return message, nil
}

func SendFile(baseUrl string, sessionId *uuid.UUID, fileInfo *FileInfo) (*Message, error) {
	finalUrl := fmt.Sprintf("%s/sessions/%s/send", baseUrl, sessionId.String())

	fileInfoJson, err := json.Marshal(*fileInfo)
	if err != nil { return nil, err }

	fileInfoJsonReader := bytes.NewReader(fileInfoJson)

	res, err := http.Post(finalUrl, "application/json", fileInfoJsonReader)
	if err != nil { return nil, err }

	message, err := unmarshalResponse(res)
	if err != nil { return nil, err }

	return message, nil
}

func ReceiveFile(baseUrl string, sessionId* uuid.UUID, fileInfo *FileInfo) (*Message, error) {
	finalUrl := fmt.Sprintf("%s/sessions/%s/receive", baseUrl, sessionId.String())

	fileInfoJson, err := json.Marshal(*fileInfo)
	if err != nil { return nil, err }

	fileInfoJsonReader := bytes.NewReader(fileInfoJson)

	res, err := http.Post(finalUrl, "application/json", fileInfoJsonReader)
	if err != nil { return nil, err }

	message, err := unmarshalResponse(res)
	if err != nil { return nil, err }

	return message, nil
}
