package testutil

import (
	"io"
)

func StubWriteCloser() io.WriteCloser {
	return &NullWriteCloser{}
}

type NullWriteCloser struct{}

func (*NullWriteCloser) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (*NullWriteCloser) Close() error {
	return nil
}
