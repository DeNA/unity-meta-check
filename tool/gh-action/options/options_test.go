package options

import (
	"encoding/json"
	"github.com/DeNA/unity-meta-check/tool/gh-action/inputs"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/testutil"
	"github.com/DeNA/unity-meta-check/util/typedpath"
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
			Env: map[string]string{
				"GITHUB_TOKEN":      "T0K3N",
				"GITHUB_WORKSPACE":  string(typedpath.NewRootRawPath("github", "workspace")),
				"GITHUB_EVENT_PATH": string(typedpath.NewRootRawPath("github", "workflows", "event.json")),
				"GITHUB_API_URL":    "https://api.github.com",
			},
			Expected: &Options{
				Inputs: in,
				Env: inputs.ActionEnv{
					GitHubToken: "T0K3N",
					Workspace:   typedpath.NewRootRawPath("github", "workspace"),
					EventPath:   typedpath.NewRootRawPath("github", "workflows", "event.json"),
					APIURL:      "https://api.github.com",
				},
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
