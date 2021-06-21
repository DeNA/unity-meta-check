package ghactions

import (
	"bytes"
	"github.com/DeNA/unity-meta-check/util/cli/opt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestWriteTo(t *testing.T) {
	cases := map[string]struct {
		Options  []opt.Option
		Expected string
	}{
		"empty (boundary)": {
			Options: []opt.Option{},
			Expected: ``,
		},
		"required string option (easiest case)": {
			Options: []opt.Option{
				opt.NewRequiredStringOption("req-str", "required string option"),
			},
			Expected: `
  "req-str":
    description: "required string option"
    required: true
`[1:],
		},
		"optional string option (easiest case)": {
			Options: []opt.Option{
				opt.NewOptionalStringOption("opt-str", "optional string option", "DEFAULT"),
			},
			Expected: `
  "opt-str":
    description: "optional string option"
    required: false
    default: "DEFAULT"
`[1:],
		},
		"required bool option (easiest case)": {
			Options: []opt.Option{
				opt.NewRequiredBoolOption("req-bool", "required bool option"),
			},
			Expected: `
  "req-bool":
    description: "required bool option"
    required: true
`[1:],
		},
		"optional bool option (easiest case)": {
			Options: []opt.Option{
				opt.NewOptionalBoolOption("opt-bool", "optional bool option", true),
			},
			Expected: `
  "opt-bool":
    description: "optional bool option"
    required: false
    default: true
`[1:],
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
