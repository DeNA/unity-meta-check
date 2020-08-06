package options

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestBuild(t *testing.T) {
	cwd, err := typedpath.Getwd()
	if err != nil {
		panic(err.Error())
	}

	cases := []struct {
		Args     []string
		Expected *Options
	}{
		{
			Args: []string{"-version"},
			Expected: &Options{
				Version: true,
			},
		},
		{
			Args: []string{"-fix-missing", "path/to/allow/autofix"},
			Expected: &Options{
				Version:      false,
				LogLevel:     logging.SeverityInfo,
				DryRun:       false,
				FixMissing:   true,
				FixDangling:  false,
				AllowedGlobs: []globs.Glob{"path/to/allow/autofix"},
				RootDirAbs:   cwd,
			},
		},
		{
			Args: []string{"-fix-dangling", "path/to/allow/autofix"},
			Expected: &Options{
				Version:      false,
				LogLevel:     logging.SeverityInfo,
				DryRun:       false,
				FixMissing:   false,
				FixDangling:  true,
				AllowedGlobs: []globs.Glob{"path/to/allow/autofix"},
				RootDirAbs:   cwd,
			},
		},
		{
			Args: []string{"-fix-missing", "-fix-dangling", "path/to/allow/autofix"},
			Expected: &Options{
				Version:      false,
				LogLevel:     logging.SeverityInfo,
				DryRun:       false,
				FixMissing:   true,
				FixDangling:  true,
				AllowedGlobs: []globs.Glob{"path/to/allow/autofix"},
				RootDirAbs:   cwd,
			},
		},
		{
			Args: []string{"-dry-run", "-fix-missing", "path/to/allow/autofix"},
			Expected: &Options{
				Version:      false,
				LogLevel:     logging.SeverityInfo,
				DryRun:       true,
				FixMissing:   true,
				FixDangling:  false,
				AllowedGlobs: []globs.Glob{"path/to/allow/autofix"},
				RootDirAbs:   cwd,
			},
		},
		{
			Args: []string{"-fix-missing", "-root-dir", string(cwd.Dir()), "path/to/allow/autofix"},
			Expected: &Options{
				Version:      false,
				LogLevel:     logging.SeverityInfo,
				DryRun:       false,
				FixMissing:   true,
				FixDangling:  false,
				AllowedGlobs: []globs.Glob{"path/to/allow/autofix"},
				RootDirAbs:   cwd.Dir(),
			},
		},
		{
			Args: []string{"-fix-missing", "-debug", "path/to/allow/autofix/"},
			Expected: &Options{
				Version:      false,
				LogLevel:     logging.SeverityDebug,
				DryRun:       false,
				FixMissing:   true,
				FixDangling:  false,
				AllowedGlobs: []globs.Glob{"path/to/allow/autofix"},
				RootDirAbs:   cwd,
			},
		},
		{
			Args: []string{"-fix-missing", "-silent", "path/to/allow/autofix/"},
			Expected: &Options{
				Version:      false,
				LogLevel:     logging.SeverityWarn,
				DryRun:       false,
				FixMissing:   true,
				FixDangling:  false,
				AllowedGlobs: []globs.Glob{"path/to/allow/autofix"},
				RootDirAbs:   cwd,
			},
		},
		{
			Args: []string{"-fix-missing", "-silent", "path/to/allow/autofix/"},
			Expected: &Options{
				Version:      false,
				LogLevel:     logging.SeverityWarn,
				DryRun:       false,
				FixMissing:   true,
				FixDangling:  false,
				AllowedGlobs: []globs.Glob{"path/to/allow/autofix"},
				RootDirAbs:   cwd,
			},
		},
		{
			Args: []string{"-fix-missing", "-debug", "-silent", "/path/to/allow/autofix"},
			Expected: &Options{
				Version:      false,
				LogLevel:     logging.SeverityDebug,
				DryRun:       false,
				FixMissing:   true,
				FixDangling:  false,
				AllowedGlobs: []globs.Glob{"path/to/allow/autofix"},
				RootDirAbs:   cwd,
			},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.Args), func(t *testing.T) {
			opts, err := Build(c.Args, cli.AnyProcInout())
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
