package files

import (
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type FileRoom struct {
	UUID       uuid.UUID
	ExpiresAt  time.Time
	FilePaths  []string
	tokens     []string
	password   string
	deleteOnce sync.Once
}

func NewFileRoom(password string, lifespanSec int) *FileRoom {
	fileRoom := &FileRoom{
		UUID:      uuid.New(),
		password:  password,
		FilePaths: make([]string, 0),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(lifespanSec)),
	}
	return fileRoom
}

func (fR *FileRoom) Delete() {
	fR.deleteOnce.Do(func() {
		for _, path := range fR.FilePaths {
			os.Remove(path)
		}
		fR.FilePaths = nil
		fR.tokens = nil
	})
}
