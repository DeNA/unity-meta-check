package git

import (
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/testutil"
	"testing"
)

func TestGlobalConfig(t *testing.T) {
	spy := testutil.SpyWriteCloser(nil)
	spyLogger := logging.SpyLogger()

	globalConfig := NewGlobalConfig(spyLogger)

	if err := globalConfig(spy, "--list"); err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	if !spy.IsClosed {
		t.Error("want true, but false")
		return
	}
}
