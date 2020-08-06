package git

import (
	"github.com/DeNA/unity-meta-check/util/logging"
	"path/filepath"
	"testing"
)

func TestNewRevParse(t *testing.T) {
	spyLogger := logging.SpyLogger()

	revParse := NewRevParse(spyLogger)

	actual, err := revParse(".", "--show-toplevel")

	if err != nil {
		t.Log(spyLogger.Logs.String())
		t.Errorf("want nil, got %#v", err)
		return
	}

	if !filepath.IsAbs(actual) {
		t.Log(spyLogger.Logs.String())
		t.Errorf("want an absolute file path, got %s", actual)
		return
	}
}
