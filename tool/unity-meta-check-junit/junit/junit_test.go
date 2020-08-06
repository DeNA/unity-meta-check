package junit

import (
	"bytes"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"testing"
)

func TestWrite(t *testing.T) {
	result := checker.NewCheckResult(
		[]typedpath.SlashPath{
			typedpath.NewSlashPathUnsafe("path/to/missing.meta"),
		},
		[]typedpath.SlashPath{
			typedpath.NewSlashPathUnsafe("path/to/dangling.meta"),
		},
	)

	buf := &bytes.Buffer{}
	if err := Write(result, 0, buf); err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	if len(buf.String()) == 0 {
		t.Error("want not empty string, got empty string")
		return
	}
}
