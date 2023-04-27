package main

import (
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/testutil"
	"github.com/google/go-cmp/cmp"
	"path/filepath"
	"testing"
)

func TestValid(t *testing.T) {
	main := NewMain()
	stdout := testutil.SpyWriteCloser(nil)
	stderr := testutil.SpyWriteCloser(nil)
	procInout := cli.ProcessInout{
		Stdin:  nil,
		Stdout: stdout,
		Stderr: stderr,
	}

	projectPath := filepath.Join("testdata", "ValidProject")
	actualExitStatus := main([]string{"-debug", projectPath}, procInout, cli.AnyEnv())

	expectedExitStatus := cli.ExitNormal
	if actualExitStatus != expectedExitStatus {
		t.Logf("stdout:\n%s", stdout.Captured.String())
		t.Logf("stderr:\n%s", stderr.Captured.String())
		t.Error(cmp.Diff(expectedExitStatus, actualExitStatus))
		return
	}

	actualStdout := stdout.Captured.String()
	expectedStdout := ""

	if actualStdout != expectedStdout {
		t.Logf("stderr:\n%s", stderr.Captured.String())
		t.Error(cmp.Diff(expectedStdout, actualStdout))
		return
	}
}

func TestInvalid(t *testing.T) {
	main := NewMain()
	stdout := testutil.SpyWriteCloser(nil)
	stderr := testutil.SpyWriteCloser(nil)
	procInout := cli.ProcessInout{
		Stdin:  nil,
		Stdout: stdout,
		Stderr: stderr,
	}

	projectPath := filepath.Join("testdata", "InvalidProject")
	actualExitStatus := main([]string{"-debug", projectPath}, procInout, cli.AnyEnv())

	expectedExitStatus := cli.ExitAbnormal
	if actualExitStatus != expectedExitStatus {
		t.Logf("stdout:\n%s", stdout.Captured.String())
		t.Logf("stderr:\n%s", stderr.Captured.String())
		t.Error(cmp.Diff(expectedExitStatus, actualExitStatus))
		return
	}

	actualStdout := stdout.Captured.String()
	expectedStdout := `missing Assets/AssetsMissing.meta
missing Assets/SubDir/SubDirFile.meta
missing LocalPackages/com.example.local.pkg/LocalPkgMissing.meta
missing Packages/com.example.pkg/PkgMissing.meta
dangling Assets/AssetsDangling.meta
dangling Assets/Dangling/Dangling.meta
dangling LocalPackages/com.example.local.pkg/LocalPkgDangling.meta
dangling Packages/com.example.pkg/PkgDangling.meta
`

	if actualStdout != expectedStdout {
		t.Logf("stderr:\n%s", stderr.Captured.String())
		t.Error(cmp.Diff(expectedStdout, actualStdout))
		return
	}
}

func TestVersion(t *testing.T) {
	main := NewMain()
	procInout := cli.AnyProcInout()

	status := main([]string{"-version"}, procInout, cli.AnyEnv())

	if status != cli.ExitNormal {
		t.Errorf("want %d, got %d", cli.ExitNormal, status)
		return
	}
}

func TestHelp(t *testing.T) {
	main := NewMain()
	procInout := cli.AnyProcInout()

	status := main([]string{"-help"}, procInout, cli.AnyEnv())

	if status != cli.ExitAbnormal {
		t.Errorf("want %d, got %d", cli.ExitNormal, status)
		return
	}
}
