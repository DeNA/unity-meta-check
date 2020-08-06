package l10n

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"strings"
	"testing"
)

func TestReadL10n(t *testing.T) {
	json := `{
	"header_status": "STATUS_HEADER",
	"header_file_path": "FILE_PATH_HEADER",
	"status_missing": "STATUS_MISSING",
	"status_dangling": "STATUS_DANGLING"
}`
	actual, err := ReadTemplate(strings.NewReader(json))

	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	expected := &Template{
		StatusHeader:   "STATUS_HEADER",
		FilePathHeader: "FILE_PATH_HEADER",
		StatusMissing:  "STATUS_MISSING",
		StatusDangling: "STATUS_DANGLING",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Error(cmp.Diff(expected, actual))
		return
	}
}


func TestWriteTemplateExample(t *testing.T) {
	buf := &bytes.Buffer{}

	WriteTemplateExample(buf)

	if buf.Len() == 0 {
		t.Error("want greater than 0, but 0")
		return
	}
}
