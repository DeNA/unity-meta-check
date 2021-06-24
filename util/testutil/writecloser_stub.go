package testutil

import (
	"io"
)

type WriteCloserStub struct {
	WriteErr error
	CloseErr error
}

func (w WriteCloserStub) Write(p []byte) (n int, err error) {
	if w.WriteErr == nil {
		return len(p), nil
	}
	return 0, w.WriteErr
}

func (w WriteCloserStub) Close() error {
	return w.CloseErr
}

func StubWriteCloser(writeErr, closeErr error) io.WriteCloser {
	return &WriteCloserStub{
		WriteErr: writeErr,
		CloseErr: closeErr,
	}
}

type NullWriteCloser struct{}

func (*NullWriteCloser) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (*NullWriteCloser) Close() error {
	return nil
}
