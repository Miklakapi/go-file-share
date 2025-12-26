package filestore

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
)

type DiskStore struct{}

var _ ports.FileStore = (*DiskStore)(nil)

func (DiskStore) Save(ctx context.Context, uploadDir, name string, r io.Reader) (path string, size int64, err error) {
	if err := ctx.Err(); err != nil {
		return "", 0, err
	}

	if uploadDir == "" {
		return "", 0, ports.ErrEmptyUploadDir
	}

	if name == "" {
		return "", 0, ports.ErrEmptyFilename
	}

	if r == nil {
		return "", 0, ports.ErrNilReader
	}

	safeName := filepath.Base(name)
	path = filepath.Join(uploadDir, safeName)

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", 0, err
	}

	dst, err := os.Create(path)
	if err != nil {
		return "", 0, err
	}
	defer func() {
		_ = dst.Close()
	}()

	size, err = io.Copy(dst, r)
	if err != nil {
		_ = os.Remove(path)
		return "", 0, err
	}

	return path, size, nil
}

func (DiskStore) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if path == "" {
		return nil, os.ErrNotExist
	}

	return os.Open(path)
}

func (DiskStore) Delete(ctx context.Context, path string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if path == "" {
		return nil
	}

	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
