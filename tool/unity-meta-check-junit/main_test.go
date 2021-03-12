package main

import (
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/testutil"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValid(t *testing.T) {
	stdout := testutil.SpyWriteCloser(nil)
	stderr := testutil.SpyWriteCloser(nil)

	procInout := cli.ProcessInout{
		Stdin:  strings.NewReader(""),
		Stdout: stdout,
		Stderr: stderr,
	}

	tmpDir, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	main := NewMain()
	actual := main([]string{filepath.Join(tmpDir, "valid.xml")}, procInout, cli.AnyEnv())

	expected := cli.ExitNormal
	if actual != expected {
		t.Logf("stdout:\n%s", stdout.Captured.String())
		t.Logf("stderr:\n%s", stderr.Captured.String())
		t.Errorf("want %#v, got %#v", expected, actual)
		return
	}
}

func TestInvalid(t *testing.T) {
	stdout := testutil.SpyWriteCloser(nil)
	stderr := testutil.SpyWriteCloser(nil)

	procInout := cli.ProcessInout{
		Stdin:  strings.NewReader(`missing path/to/missing.meta
dangling path/to/dangling.meta`),
		Stdout: stdout,
		Stderr: stderr,
	}

	tmpDir, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	main := NewMain()
	actual := main([]string{filepath.Join(tmpDir, "invalid.xml")}, procInout, cli.AnyEnv())

	expected := cli.ExitAbnormal
	if actual != expected {
		t.Logf("stdout:\n%s", stdout.Captured.String())
		t.Logf("stderr:\n%s", stderr.Captured.String())
		t.Errorf("want %#v, got %#v", expected, actual)
		return
	}
}

func TestVersion(t *testing.T) {
	main := NewMain()
	stdout := testutil.SpyWriteCloser(nil)
	stderr := testutil.SpyWriteCloser(nil)
	procInout := cli.ProcessInout{
		Stdin:  strings.NewReader(""),
		Stdout: stdout,
		Stderr: stderr,
	}

	actual := main([]string{"-version"}, procInout, cli.AnyEnv())

	expected := cli.ExitNormal
	if actual != expected {
		t.Logf("stdout:\n%s", stdout.Captured.String())
		t.Logf("stderr:\n%s", stderr.Captured.String())
		t.Errorf("want %#v, got %#v", expected, actual)
		return
	}
}
