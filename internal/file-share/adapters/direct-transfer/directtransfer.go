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

	dT.mu.RLock()
	_, ok := dT.connections[code]
	dT.mu.RUnlock()
	if ok {
		return nil, ports.ErrTransferCodeExists
	}

	c := connection{}

	dT.mu.Lock()
	dT.connections[code] = &c
	dT.mu.Unlock()
	panic("TODO")
}

func (dT *DirectTransfer) Send(code string, filename string, src io.Reader) error {
	if !codeOk(code) {
		return ports.ErrTransferCodeInvalidLength
	}

	dT.mu.RLock()
	connection, ok := dT.connections[code]
	dT.mu.RUnlock()

	if !ok {
		return ports.ErrTransferCodeNotFound
	}

	pr, pw := io.Pipe()
	defer func() { _ = pw.Close() }()
	transfer := ports.Transfer{Reader: pr, Filename: filename}

	dT.mu.Lock()
	connection.transfer = transfer
	dT.mu.Unlock()

	_, err := io.Copy(pw, src)

	return err
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
