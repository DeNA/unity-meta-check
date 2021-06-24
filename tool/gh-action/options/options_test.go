package options

import (
	"encoding/json"
	"github.com/DeNA/unity-meta-check/tool/gh-action/inputs"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/testutil"
	"github.com/google/go-cmp/cmp"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	in := inputs.Inputs{}
	inputsJson, err := json.Marshal(in)
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	cases := map[string]struct {
		Args     []string
		Env      map[string]string
		Expected *Options
	}{
		"easiest case": {
			Args: []string{"-inputs-json", string(inputsJson)},
			Env:  map[string]string{"GITHUB_TOKEN": "T0K3N"},
			Expected: &Options{
				UnsafeInputs: in,
				Token:        "T0K3N",
			},
		},
		"-version": {
			Args: []string{"-version"},
			Expected: &Options{
				Version: true,
			},
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

			parse := NewParser()
			opts, err := parse(c.Args, procInout, cli.StubEnv(c.Env))
			if err != nil {
				t.Errorf("want nil, got %#v", err)
				return
			}

			if !reflect.DeepEqual(opts, c.Expected) {
				t.Error(cmp.Diff(c.Expected, opts))
			}
		})
	}
}
