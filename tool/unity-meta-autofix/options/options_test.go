package options

import (
	"github.com/DeNA/unity-meta-check/options"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestBuild(t *testing.T) {
	cases := map[string]struct {
		Args     []string
		Cwd      typedpath.SlashPath
		Expected *Options
	}{
		"-version": {
			Args: []string{"-version"},
			Expected: &Options{
				Version: true,
			},
		},
		"only 1 glob": {
			Args: []string{string(typedpath.NewRawPath("path", "to", "allow", "autofix"))},
			Cwd:  "/abs",
			Expected: &Options{
				Version:      false,
				LogLevel:     logging.SeverityInfo,
				DryRun:       false,
				AllowedGlobs: []globs.Glob{"path/to/allow/autofix"},
				RootDirAbs:   typedpath.NewRootRawPath("abs"),
			},
		},
		"-dry-run": {
			Args: []string{"-dry-run", string(typedpath.NewRawPath("path", "to", "allow", "autofix"))},
			Cwd:  "/abs",
			Expected: &Options{
				Version:      false,
				LogLevel:     logging.SeverityInfo,
				DryRun:       true,
				AllowedGlobs: []globs.Glob{"path/to/allow/autofix"},
				RootDirAbs:   typedpath.NewRootRawPath("abs"),
			},
		},
		"-root-dir": {
			Args: []string{"-root-dir", string(typedpath.NewRootRawPath("root", "dir")), string(typedpath.NewRawPath("path", "to", "allow", "autofix"))},
			Expected: &Options{
				Version:      false,
				LogLevel:     logging.SeverityInfo,
				DryRun:       false,
				AllowedGlobs: []globs.Glob{"path/to/allow/autofix"},
				RootDirAbs:   typedpath.NewRootRawPath("root", "dir"),
			},
		},
		"-debug": {
			Args: []string{"-debug", string(typedpath.NewRawPath("path", "to", "allow", "autofix"))},
			Cwd:  "/abs",
			Expected: &Options{
				Version:      false,
				LogLevel:     logging.SeverityDebug,
				DryRun:       false,
				AllowedGlobs: []globs.Glob{"path/to/allow/autofix"},
				RootDirAbs:   typedpath.NewRootRawPath("abs"),
			},
		},
		"-silent": {
			Args: []string{"-silent", string(typedpath.NewRawPath("path", "to", "allow", "autofix"))},
			Cwd:  "/abs",
			Expected: &Options{
				Version:      false,
				LogLevel:     logging.SeverityWarn,
				DryRun:       false,
				AllowedGlobs: []globs.Glob{"path/to/allow/autofix"},
				RootDirAbs:   typedpath.NewRootRawPath("abs"),
			},
		},
		"both -debug and -silent": {
			Args: []string{"-debug", "-silent", string(typedpath.NewRawPath("path", "to", "allow", "autofix"))},
			Cwd:  "/abs",
			Expected: &Options{
				Version:      false,
				LogLevel:     logging.SeverityDebug,
				DryRun:       false,
				AllowedGlobs: []globs.Glob{"path/to/allow/autofix"},
				RootDirAbs:   typedpath.NewRootRawPath("abs"),
			},
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			parse := NewParser(options.FakeRootDirValidator(c.Cwd.ToRaw()))

			opts, err := parse(c.Args, cli.AnyProcInout())
			if err != nil {
				t.Errorf("want nil, got %#v", err)
				return
			}

			if !reflect.DeepEqual(opts, c.Expected) {
				t.Error(cmp.Diff(c.Expected, opts))
				return
			}
		})
	}
}
