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

	if err := globalConfig(spy, "--default", "", "--get", "safe.directory"); err != nil {
		t.Errorf("want nil, got %#v", err)
		t.Log(spyLogger.Logs.String())
		return
	}

	if !spy.IsClosed {
		t.Error("want true, but false")
		return
	}
}
