package directtransfer

import (
	"context"
	"io"
	"sync"

	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
)

type connection struct {
	transfer ports.Transfer
}

type DirectTransfer struct {
	mu          sync.RWMutex
	connections map[string]*connection
}

var _ ports.DirectTransfer = (*DirectTransfer)(nil)

func New() *DirectTransfer {
	return &DirectTransfer{
		connections: make(map[string]*connection, 5),
	}
}

func (dT *DirectTransfer) Receive(ctx context.Context, code string) (*ports.Transfer, error) {
	if !codeOk(code) {
		return nil, ports.ErrTransferCodeInvalidLength
	}
	panic("TODO")
}

func (dT *DirectTransfer) Send(code string, filename string, src io.Reader) error {
	if !codeOk(code) {
		return ports.ErrTransferCodeInvalidLength
	}
	panic("TODO")
}

func (dT *DirectTransfer) Cancel(code string) {
	if !codeOk(code) {
		return
	}
	panic("TODO")
}

func codeOk(code string) bool {
	l := len(code)
	if l < 16 || l > 16 {
		return false
	}
	return true
}
