package file

import (
	"errors"
	"github.com/nu7hatch/gouuid"
)

type FileMap map[uuid.UUID]*File

var SendingFileMap FileMap
var ReceivingFileMap FileMap

func (fileMap FileMap) checkAndNewFileMap() {
	if fileMap == nil {
		fileMap = make(FileMap)
	}
}

func (fileMap FileMap) NewFileId() (*uuid.UUID, error) {
	fileId, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return fileId, nil
}

func (fileMap FileMap) PutFile(file *File) error {
	fileMap.checkAndNewFileMap()

	if file.FileId == nil {
		return errors.New("invalid fileId")
	}

	fileMap[*file.FileId] = file

	return nil
}

func (fileMap FileMap) NewFile() (*File, error) {
	fileId, err := fileMap.NewFileId()
	if err != nil { return nil, err }

	file := new(File)
	file.FileId = fileId
	err = fileMap.PutFile(file)
	if err != nil { return nil, err }

	return file, nil
}

func (fileMap FileMap) GetFile(fileId *uuid.UUID) (*File, error) {
	fileMap.checkAndNewFileMap()

	return fileMap[*fileId], nil
}
