package prefix

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"io"
	"testing"
)

func TestWriterEmpty(t *testing.T) {
	buf := &bytes.Buffer{}

	writer := NewWriter("PREFIX:", buf)

	if _, err := io.WriteString(writer, ""); err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	actual := buf.String()
	expected := "PREFIX:"
	if actual != expected {
		t.Error(cmp.Diff(expected, actual))
		return
	}
}

func TestWriterFirstLine(t *testing.T) {
	buf := &bytes.Buffer{}

	writer := NewWriter("PREFIX:", buf)

	if _, err := io.WriteString(writer, "LINE1"); err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	actual := buf.String()
	expected := "PREFIX:LINE1"
	if actual != expected {
		t.Error(cmp.Diff(expected, actual))
		return
	}
}

func TestWriterFirstLineEnd(t *testing.T) {
	buf := &bytes.Buffer{}

	writer := NewWriter("PREFIX:", buf)

	if _, err := io.WriteString(writer, "LINE1\n"); err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	actual := buf.String()
	expected := "PREFIX:LINE1\nPREFIX:"
	if actual != expected {
		t.Error(cmp.Diff(expected, actual))
		return
	}
}

func TestWriterFirstLine2(t *testing.T) {
	buf := &bytes.Buffer{}

	writer := NewWriter("PREFIX:", buf)

	if _, err := io.WriteString(writer, "LINE1\nLINE2"); err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	actual := buf.String()
	expected := "PREFIX:LINE1\nPREFIX:LINE2"
	if actual != expected {
		t.Error(cmp.Diff(expected, actual))
		return
	}
}
