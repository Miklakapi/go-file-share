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

func (DiskStore) ClearAll(ctx context.Context, uploadDir string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if uploadDir == "" {
		return nil
	}

	info, err := os.Stat(uploadDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return ports.ErrInvalidUploadDir
	}

	entries, err := os.ReadDir(uploadDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return err
		}

		path := filepath.Join(uploadDir, entry.Name())

		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	return nil
}

func (DiskStore) Save(ctx context.Context, uploadDir, name string, r io.Reader) (string, int64, error) {
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
	path := filepath.Join(uploadDir, safeName)

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

	size, err := io.Copy(dst, r)
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

func (DiskStore) Exists(ctx context.Context, path string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	if path == "" {
		return false, os.ErrNotExist
	}

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
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
