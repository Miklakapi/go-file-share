package directtransfer

import (
	"context"
	"io"
	"sync"

	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
)

type connection struct {
	transfer    chan *ports.Transfer
	sessionDone chan struct{}
	once        sync.Once
	mu          sync.Mutex
	pw          *io.PipeWriter
}

func newConnection() *connection {
	return &connection{
		transfer:    make(chan *ports.Transfer, 1),
		sessionDone: make(chan struct{}),
	}
}

func (c *connection) close(err error) {
	c.once.Do(func() {
		c.mu.Lock()
		if c.pw != nil {
			if err != nil {
				_ = c.pw.CloseWithError(err)
			} else {
				_ = c.pw.Close()
			}
			c.pw = nil
		}
		c.mu.Unlock()

		close(c.sessionDone)
	})
}

func (c *connection) setWriter(pw *io.PipeWriter) {
	c.mu.Lock()
	c.pw = pw
	c.mu.Unlock()
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

	c := newConnection()

	dT.mu.Lock()
	dT.connections[code] = c
	dT.mu.Unlock()

	go func() {
		<-c.sessionDone
		dT.mu.Lock()
		if cur, ok := dT.connections[code]; ok && cur == c {
			delete(dT.connections, code)
		}
		dT.mu.Unlock()
	}()

	select {
	case tr := <-c.transfer:
		return tr, nil

	case <-c.sessionDone:
		return nil, context.Canceled

	case <-ctx.Done():
		c.close(ctx.Err())
		return nil, ctx.Err()
	}
}

func (dT *DirectTransfer) Send(ctx context.Context, code string, filename string, src io.Reader) error {
	if !codeOk(code) {
		return ports.ErrTransferCodeInvalidLength
	}

	dT.mu.RLock()
	c, ok := dT.connections[code]
	dT.mu.RUnlock()
	if !ok {
		return ports.ErrTransferCodeNotFound
	}

	pr, pw := io.Pipe()
	c.setWriter(pw)
	defer func() {
		c.mu.Lock()
		if c.pw == pw {
			c.pw = nil
		}
		c.mu.Unlock()
	}()

	tr := &ports.Transfer{Reader: pr, Filename: filename}
	select {
	case c.transfer <- tr:

	case <-c.sessionDone:
		_ = pr.Close()
		_ = pw.CloseWithError(context.Canceled)
		return context.Canceled
	case <-ctx.Done():
		_ = pr.Close()
		_ = pw.CloseWithError(ctx.Err())
		c.close(ctx.Err())
		return ctx.Err()
	}

	copyDone := make(chan error, 1)
	go func() {
		_, err := io.Copy(pw, src)
		if err != nil {
			_ = pw.CloseWithError(err)
		} else {
			_ = pw.Close()
		}
		copyDone <- err
	}()

	select {
	case err := <-copyDone:
		if err != nil {
			c.close(err)
			return err
		}
		c.close(nil)
		return nil

	case <-c.sessionDone:
		_ = pw.CloseWithError(context.Canceled)
		err := <-copyDone
		if err == nil {
			err = context.Canceled
		}
		c.close(err)
		return context.Canceled

	case <-ctx.Done():
		_ = pw.CloseWithError(ctx.Err())
		err := <-copyDone
		c.close(ctx.Err())
		if err == nil {
			return ctx.Err()
		}
		return ctx.Err()
	}
}

func (dT *DirectTransfer) Cancel(code string) error {
	if !codeOk(code) {
		return ports.ErrTransferCodeInvalidLength
	}

	dT.mu.RLock()
	c, ok := dT.connections[code]
	dT.mu.RUnlock()

	if !ok {
		return ports.ErrTransferCodeNotFound
	}

	c.close(context.Canceled)
	return nil
}

func codeOk(code string) bool {
	return len(code) == 16
}
