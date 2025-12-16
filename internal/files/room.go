package files

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type FileRoom struct {
	Id         uuid.UUID
	password   string ``
	FilePaths  []string
	Lifespan   *time.Ticker
	FileHub    *FileHub
	AddFile    chan []*multipart.FileHeader
	RemoveFile chan []*multipart.FileHeader
}

func NewFileContainer(password string, lifespan int, fileHub *FileHub) *FileRoom {
	fileContainre := &FileRoom{
		Id:         uuid.New(),
		password:   password,
		FilePaths:  make([]string, 0),
		Lifespan:   time.NewTicker(time.Duration(lifespan)),
		FileHub:    fileHub,
		AddFile:    make(chan []*multipart.FileHeader),
		RemoveFile: make(chan []*multipart.FileHeader),
	}
	go fileContainre.run()
	return fileContainre
}

func (fC *FileRoom) run() {
	defer func() {
		fC.Lifespan.Stop()
		fC.FileHub.Unregister <- fC

	}()
	<-fC.Lifespan.C
}
