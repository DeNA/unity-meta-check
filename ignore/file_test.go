package ignore

import (
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	ignoreFile := `
# this is a comment
foo
bar/path # also this is a comment
baz/ # this will be trimmed
/qux # also this will be trimmed
`

	reader := strings.NewReader(ignoreFile)
	actual, err := Read(reader)

	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	expected := []globs.Glob{"foo", "bar/path", "baz", "qux"}
	if !reflect.DeepEqual(actual, expected) {
		t.Error(cmp.Diff(expected, actual))
		return
	}
}
