package git

import (
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/testutil"
	"testing"
)

func TestLsFiles(t *testing.T) {
	spy := testutil.SpyWriteCloser(nil)
	spyLogger := logging.SpyLogger()

	lsFiles := NewLsFiles(spyLogger)

	if err := lsFiles(".", []string{}, spy); err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	stdout := spy.Captured.String()
	if len(stdout) == 0 {
		t.Log(spyLogger.Logs.String())
		t.Error("want stdout > 0, but == 0")
		return
	}
	if !spy.IsClosed {
		t.Error("want true, but false")
		return
	}
}
