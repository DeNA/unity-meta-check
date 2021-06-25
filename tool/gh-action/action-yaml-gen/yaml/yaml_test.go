package yaml

import (
	"bytes"
	"github.com/DeNA/unity-meta-check/version"
	"github.com/google/go-cmp/cmp"
	"os"
	"strings"
	"testing"
)

// NOTE: This test is fragile, but we can use like Golden Testing.
func TestRecentActionYAML(t *testing.T) {
	actualRawBytes, err := os.ReadFile("../testdata/action.yml")
	if err != nil {
		t.Errorf("want nil, got %v", err)
		return
	}

	// NOTE: On Windows, Git may replace \n to \r\n on specific configurations.
	actual := strings.NewReplacer("\r", "", "%VERSION%", version.Version).Replace(string(actualRawBytes))

	buf := &bytes.Buffer{}
	if _, err := WriteTo(buf); err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	if buf.String() != actual {
		t.Error(cmp.Diff(buf.String(), actual))
	}
}
