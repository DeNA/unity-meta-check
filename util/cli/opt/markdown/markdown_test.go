package markdown

import (
	"bytes"
	"github.com/DeNA/unity-meta-check/util/cli/opt"
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

func TestWriteTo(t *testing.T) {
	cases := map[string]struct {
		Options  []opt.Option
		Expected string
	}{
		"empty (boundary case)": {
			Options: []opt.Option{},
			Expected: strings.Join([]string{
				"| Option | Description | Required or Default Value |\n",
				"|:-------|:------------|:--------------------------|\n",
			}, ""),
		},
		"optional string option (easiest case)": {
			Options: []opt.Option{
				opt.NewOptionalStringOption("opt-str", "optional string option", "DEFAULT"),
			},
			Expected: strings.Join([]string{
				"| Option | Description | Required or Default Value |\n",
				"|:-------|:------------|:--------------------------|\n",
				"| `--opt-str <string>` | optional string option | optional (default: `\"DEFAULT\"`) |\n",
			}, ""),
		},
		"required string option (easiest case)": {
			Options: []opt.Option{
				opt.NewRequiredStringOption("req-str", "required string option"),
			},
			Expected: strings.Join([]string{
				"| Option | Description | Required or Default Value |\n",
				"|:-------|:------------|:--------------------------|\n",
				"| `--req-str <string>` | required string option | required |\n",
			}, ""),
		},
		"optional bool option (easiest case)": {
			Options: []opt.Option{
				opt.NewOptionalBoolOption("opt-bool", "optional bool option", true),
			},
			Expected: strings.Join([]string{
				"| Option | Description | Required or Default Value |\n",
				"|:-------|:------------|:--------------------------|\n",
				"| `--opt-bool` | optional bool option | optional (default: `true`) |\n",
			}, ""),
		},
		"required bool option (easiest case)": {
			Options: []opt.Option{
				opt.NewRequiredBoolOption("req-bool", "required bool option"),
			},
			Expected: strings.Join([]string{
				"| Option | Description | Required or Default Value |\n",
				"|:-------|:------------|:--------------------------|\n",
				"| `--req-bool` | required bool option | required |\n",
			}, ""),
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			buf := &bytes.Buffer{}

			_, err := WriteTo(buf, c.Options...)
			if err != nil {
				t.Errorf("want nil, got %#v", err)
				return
			}

			if buf.String() != c.Expected {
				t.Error(cmp.Diff(c.Expected, buf.String()))
			}
		})
	}
}
