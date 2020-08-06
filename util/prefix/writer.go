package prefix

import (
	"bytes"
	"io"
)

type Writer struct {
	prefix    string
	base      io.Writer
	buf       *bytes.Buffer
	firstDone bool
}

func NewWriter(prefix string, base io.Writer) io.Writer {
	return &Writer{prefix, base, &bytes.Buffer{}, false}
}

func (w *Writer) Write(bs []byte) (int, error) {
	w.buf.Reset()

	if !w.firstDone {
		w.buf.WriteString(w.prefix)
		w.firstDone = true
	}

	for _, b := range bs {
		if b == '\n' {
			_ = w.buf.WriteByte(b)
			w.buf.WriteString(w.prefix)
		} else {
			w.buf.WriteByte(b)
		}
	}

	n, err := w.buf.WriteTo(w.base)
	return int(n), err
}
