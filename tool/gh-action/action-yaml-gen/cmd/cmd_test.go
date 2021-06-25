package cmd

import (
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/testutil"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestNewMain(t *testing.T) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "action-yaml-gen-test.*")
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	stdin := io.NopCloser(strings.NewReader(""))
	stdout := testutil.SpyWriteCloser(nil)
	stderr := testutil.SpyWriteCloser(nil)
	procInput := cli.ProcessInout{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}
	env := cli.StubEnv(map[string]string{})

	actual := Main([]string{path.Join(tmpDir, "action.yml")}, procInput, env)

	if actual != cli.ExitNormal {
		t.Log(stdout.Captured.String())
		t.Log(stderr.Captured.String())
		t.Errorf("want %d, got %d", cli.ExitAbnormal, actual)
	}
}
