package file

import (
	"errors"
	"github.com/nu7hatch/gouuid"
)

var FileMap map[uuid.UUID]*File

func checkAndNewFileMap() {
	if FileMap == nil {
		FileMap = make(map[uuid.UUID]*File)
	}
}

func NewFileId() (*uuid.UUID, error) {
	fileId, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return fileId, nil
}

func PutFile(file *File) error {
	checkAndNewFileMap()

	if file.FileId == nil {
		return errors.New("invalid fileId")
	}

	FileMap[*file.FileId] = file

	return nil
}

func NewFile() (*File, error) {
	fileId, err := NewFileId()
	if err != nil { return nil, err }

	file := new(File)
	file.FileId = fileId
	err = PutFile(file)
	if err != nil { return nil, err }

	return file, nil
}

func GetFile(fileId *uuid.UUID) (*File, error) {
	checkAndNewFileMap()

	return FileMap[*fileId], nil
}
