package cmd

import (
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/testutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidDryRun(t *testing.T) {
	main := NewMain()
	stdout := testutil.SpyWriteCloser(nil)
	stderr := testutil.SpyWriteCloser(nil)
	procInout := cli.ProcessInout{
		Stdin:  strings.NewReader(""),
		Stdout: stdout,
		Stderr: stderr,
	}

	rootDir := filepath.Join("testdata", "ValidProject")
	actual := main([]string{"-debug", "-dry-run", "-fix-missing", "-fix-dangling", "-root-dir", rootDir, "Assets/*"}, procInout, cli.AnyEnv())

	expected := cli.ExitNormal
	if actual != expected {
		t.Logf("stdout:\n%s", stdout.Captured.String())
		t.Logf("stderr:\n%s", stderr.Captured.String())
		t.Errorf("want %#v, got %#v", expected, actual)
		return
	}
}

func TestInvalidDryRun(t *testing.T) {
	main := NewMain()
	stdout := testutil.SpyWriteCloser(nil)
	stderr := testutil.SpyWriteCloser(nil)
	procInout := cli.ProcessInout{
		Stdin:  strings.NewReader(`missing Missing.meta
dangling Dangling.meta`),
		Stdout: stdout,
		Stderr: stderr,
	}

	rootDir := filepath.Join("testdata", "InvalidProject")
	actual := main([]string{"-debug", "-dry-run", "-fix-missing", "-fix-dangling", "-root-dir", rootDir, "Assets/*"}, procInout, cli.AnyEnv())

	expected := cli.ExitNormal
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
