package cmd

import (
	"encoding/json"
	"github.com/DeNA/unity-meta-check/tool/gh-action/inputs"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/testutil"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"io"
	"strings"
	"testing"
)

func TestNewMain(t *testing.T) {
	cases := map[string]struct {
		Args     []string
		Env      map[string]string
		Expected cli.ExitStatus
	}{
		"check-only auto-detect ValidProject": {
			Args: []string{
				"-inputs-json", string(MustMarshal(inputs.Inputs{
					LogLevel:   "INFO",
					TargetType: "auto-detect",
					TargetPath: typedpath.NewRawPath("testdata", "ValidProject"),
				})),
			},
			Env:      map[string]string{"GITHUB_TOKEN": "T0K3N"},
			Expected: cli.ExitNormal,
		},
		"check-only auto-detect InvalidProject": {
			Args: []string{
				"-inputs-json", string(MustMarshal(inputs.Inputs{
					LogLevel:   "DEBUG",
					TargetType: "auto-detect",
					TargetPath: typedpath.NewRawPath("testdata", "InvalidProject"),
				})),
			},
			Env:      map[string]string{"GITHUB_TOKEN": "T0K3N"},
			Expected: cli.ExitAbnormal,
		},
		"check-only auto-detect ValidSubDir": {
			Args: []string{
				"-inputs-json", string(MustMarshal(inputs.Inputs{
					LogLevel:   "INFO",
					TargetType: "auto-detect",
					TargetPath: typedpath.NewRawPath("testdata", "ValidProject", "LocalPackages", "com.example.local.pkg"),
				})),
			},
			Env:      map[string]string{"GITHUB_TOKEN": "T0K3N"},
			Expected: cli.ExitNormal,
		},
		"check-only auto-detect InvalidSubDir": {
			Args: []string{
				"-inputs-json", string(MustMarshal(inputs.Inputs{
					LogLevel:   "INFO",
					TargetType: "auto-detect",
					TargetPath: typedpath.NewRawPath("testdata", "InvalidProject", "LocalPackages", "com.example.local.pkg"),
				})),
			},
			Env:      map[string]string{"GITHUB_TOKEN": "T0K3N"},
			Expected: cli.ExitAbnormal,
		},
		"-version": {
			Args:     []string{"-version"},
			Env:      map[string]string{"GITHUB_TOKEN": "T0K3N"},
			Expected: cli.ExitNormal,
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			stdin := io.NopCloser(strings.NewReader(""))
			stdout := testutil.SpyWriteCloser(nil)
			stderr := testutil.SpyWriteCloser(nil)
			procInout := cli.ProcessInout{
				Stdin:  stdin,
				Stdout: stdout,
				Stderr: stderr,
			}
			env := cli.StubEnv(c.Env)

			actual := Main(c.Args, procInout, env)
			if actual != c.Expected {
				t.Log(stdout.Captured.String())
				t.Log(stderr.Captured.String())
				t.Errorf("want %d, got %d", c.Expected, actual)
			}
		})
	}
}

func MustMarshal(o interface{}) []byte {
	b, err := json.Marshal(o)
	if err != nil {
		panic(err.Error())
	}
	return b
}
