package yaml

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"os"
	"reflect"
	"testing"
)

// NOTE: This test is fragile, but we can use like Golden Testing.
func TestRecentActionYAML(t *testing.T) {
	actual, err := os.ReadFile("../testdata/action.yml")
	if err != nil {
		t.Errorf("want nil, got %v", err)
		return
	}

	buf := &bytes.Buffer{}
	if _, err := WriteTo(buf); err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	if !reflect.DeepEqual(buf.Bytes(), actual) {
		t.Error(cmp.Diff(buf.String(), string(actual)))
	}
}
