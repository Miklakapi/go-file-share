package filestore

import (
	"context"
	"fmt"
	"io"

	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
)

type TestStore struct{}

var _ ports.FileStore = (*TestStore)(nil)

func (TestStore) Save(ctx context.Context, uploadDir, name string, r io.Reader) (path string, size int64, err error) {
	fmt.Printf("Saving file")
	return name, 0, nil
}

func (TestStore) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	fmt.Printf("Opening file")
	return nil, nil
}

func (TestStore) Delete(ctx context.Context, path string) error {
	fmt.Printf("Deleting file")
	return nil
}
