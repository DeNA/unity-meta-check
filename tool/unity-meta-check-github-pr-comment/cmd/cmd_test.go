package cmd

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/testutil"
	"strings"
	"testing"
)

func TestValid(t *testing.T) {
	testEnv, err := testutil.GetTestEnv()
	if err != nil {
		t.Log(err.Error())
		t.Skip("no environment variables for tests")
		return
	}

	stdout := testutil.SpyWriteCloser(nil)
	stderr := testutil.SpyWriteCloser(nil)

	procInout := cli.ProcessInout{
		Stdin:  strings.NewReader(""),
		Stdout: stdout,
		Stderr: stderr,
	}

	main := NewMain()
	actual := main([]string{
		"-debug",
		"-owner", testEnv.Owner,
		"-repo", testEnv.Repo,
		"-pull", testEnv.Pull,
		"-api-endpoint", testEnv.ApiEndpoint,
	}, procInout, func(key string) string {
		switch key {
		case "GITHUB_TOKEN":
			return testEnv.Token
		default:
			panic(fmt.Sprintf("unsupported key: %q", key))
		}
	})

	expected := cli.ExitNormal
	if actual != expected {
		t.Logf("stdout:\n%s", stdout.Captured.String())
		t.Logf("stderr:\n%s", stderr.Captured.String())
		t.Errorf("want %#v, got %#v", expected, actual)
		return
	}
}

func TestInvalid(t *testing.T) {
	testEnv, err := testutil.GetTestEnv()
	if err != nil {
		t.Log(err.Error())
		t.Skip("no environment variables for tests")
		return
	}

	stdout := testutil.SpyWriteCloser(nil)
	stderr := testutil.SpyWriteCloser(nil)

	procInout := cli.ProcessInout{
		Stdin: strings.NewReader(`missing path/to/missing.meta
dangling path/to/dangling.meta`),
		Stdout: stdout,
		Stderr: stderr,
	}

	main := NewMain()
	actual := main([]string{
		"-debug",
		"-owner", testEnv.Owner,
		"-repo", testEnv.Repo,
		"-pull", testEnv.Pull,
		"-api-endpoint", testEnv.ApiEndpoint,
	}, procInout, func(key string) string {
		switch key {
		case "GITHUB_TOKEN":
			return testEnv.Token
		default:
			panic(fmt.Sprintf("unsupported key: %q", key))
		}
	})

	expected := cli.ExitAbnormal
	if actual != expected {
		t.Logf("stdout:\n%s", stdout.Captured.String())
		t.Logf("stderr:\n%s", stderr.Captured.String())
		t.Errorf("want %#v, got %#v", expected, actual)
		return
	}
}

func TestVersion(t *testing.T) {
	stdout := testutil.SpyWriteCloser(nil)
	stderr := testutil.SpyWriteCloser(nil)
	procInout := cli.ProcessInout{
		Stdin:  strings.NewReader(""),
		Stdout: stdout,
		Stderr: stderr,
	}

	main := NewMain()
	actual := main([]string{"-version"}, procInout, cli.AnyEnv())

	expected := cli.ExitNormal
	if actual != expected {
		t.Logf("stdout:\n%s", stdout.Captured.String())
		t.Logf("stderr:\n%s", stderr.Captured.String())
		t.Errorf("want %#v, got %#v", expected, actual)
		return
	}
}
