package files

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmptyUploadDir = errors.New("upload dir is empty")
	ErrEmptyFilename  = errors.New("filename is empty")
	ErrNilReader      = errors.New("file reader is nil")
)

type FileRoomFile struct {
	ID        uuid.UUID
	Path      string
	Name      string
	Size      int64
	CreatedAt time.Time
}

func NewFileRoomFile(uploadDir string, name string, reader io.Reader) (*FileRoomFile, error) {
	uploadDir = strings.TrimSpace(uploadDir)
	if uploadDir == "" {
		return nil, ErrEmptyUploadDir
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyFilename
	}
	if reader == nil {
		return nil, ErrNilReader
	}

	id := uuid.New()
	safeName := filepath.Base(name)
	ext := filepath.Ext(safeName)
	path := filepath.Join(uploadDir, id.String()+ext)

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	dst, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = dst.Close() }()

	size, err := io.Copy(dst, reader)
	if err != nil {
		_ = os.Remove(path)
		return nil, err
	}

	return &FileRoomFile{
		ID:        id,
		Path:      path,
		Name:      safeName,
		Size:      size,
		CreatedAt: time.Now(),
	}, nil
}

func (f *FileRoomFile) Delete() error {
	if f == nil || strings.TrimSpace(f.Path) == "" {
		return nil
	}

	err := os.Remove(f.Path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
