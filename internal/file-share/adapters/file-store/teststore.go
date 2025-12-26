package filestore

import (
	"context"
	"io"

	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
)

type TestStore struct{}

var _ ports.FileStore = (*TestStore)(nil)

func (TestStore) Save(ctx context.Context, uploadDir, name string, r io.Reader) (path string, size int64, err error) {
	panic("TODO")
}

func (TestStore) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	panic("TODO")
}

func (TestStore) Delete(ctx context.Context, path string) error {
	panic("TODO")
}
