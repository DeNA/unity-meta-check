package testutil

import (
	"bytes"
	"io"
)

func SpyWriteCloser(base io.WriteCloser) *WriteCloserSpy {
	var inherited io.WriteCloser
	if base != nil {
		inherited = base
	} else {
		inherited = StubWriteCloser()
	}

	captured := &bytes.Buffer{}

	return &WriteCloserSpy{
		Inherited: inherited,
		Captured:  captured,
		IsClosed:  false,
	}
}

type WriteCloserSpy struct {
	Inherited io.WriteCloser
	Captured  *bytes.Buffer
	IsClosed  bool
}

func (s *WriteCloserSpy) Write(p []byte) (n int, err error) {
	s.Captured.Write(p)
	return s.Inherited.Write(p)
}

func (s *WriteCloserSpy) Close() error {
	s.IsClosed = true
	return s.Inherited.Close()
}
