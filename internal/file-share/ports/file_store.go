package ports

import (
	"context"
	"io"
)

type FileStore interface {
	ClearAll(ctx context.Context, uploadDir string) error
	Save(ctx context.Context, uploadDir, name string, r io.Reader) (path string, size int64, err error)
	Open(ctx context.Context, path string) (io.ReadCloser, error)
	Exists(ctx context.Context, path string) (bool, error)
	Delete(ctx context.Context, path string) error
}
