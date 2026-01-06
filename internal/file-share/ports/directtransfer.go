package ports

import (
	"context"
	"io"
)

type Transfer struct {
	Reader   io.Reader
	Filename string
}

type DirectTransfer interface {
	Receive(ctx context.Context, code string) (Transfer, error)
	Send(code string, filename string, src io.Reader) error
	Cancel(code string)
}
