package ports

import (
	"context"
	"io"
)

type FileStore interface {
	Save(ctx context.Context, uploadDir, name string, r io.Reader) (path string, size int64, err error)
	Open(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
}
