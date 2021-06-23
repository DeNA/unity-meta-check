package options

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/testutil"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"strings"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	cases := []struct {
		Args         []string
		TargetType   checker.TargetType
		RootDirAbs   typedpath.RawPath
		IgnoredPaths []globs.Glob

		Expected *Options
	}{
		{
			Args:         []string{},
			TargetType:   checker.TargetTypeIsUnityProjectRootDirectory,
			RootDirAbs:   typedpath.NewRawPath("path", "to", "unity", "proj"),
			IgnoredPaths: []globs.Glob{},
			Expected: &Options{
				Version:        false,
				LogLevel:       logging.SeverityInfo,
				TargetType:     checker.TargetTypeIsUnityProjectRootDirectory,
				IgnoreDangling: false,
				IgnoreCase:     true,
				IgnoredPaths:   []globs.Glob{},
				RootDirAbs:     typedpath.NewRawPath("path", "to", "unity", "proj"),
			},
		},
		{
			Args:         []string{},
			TargetType:   checker.TargetTypeIsUnityProjectSubDirectory,
			RootDirAbs:   typedpath.NewRawPath("path", "to", "unity", "proj", "sub", "dir"),
			IgnoredPaths: []globs.Glob{},
			Expected: &Options{
				Version:        false,
				LogLevel:       logging.SeverityInfo,
				TargetType:     checker.TargetTypeIsUnityProjectSubDirectory,
				IgnoreDangling: false,
				IgnoreCase:     true,
				IgnoredPaths:   []globs.Glob{},
				RootDirAbs:     typedpath.NewRawPath("path", "to", "unity", "proj", "sub", "dir"),
			},
		},
		{
			Args:         []string{"-version"},
			TargetType:   checker.TargetTypeIsUnityProjectRootDirectory,
			RootDirAbs:   typedpath.NewRawPath("path", "to", "unity", "proj"),
			IgnoredPaths: []globs.Glob{},
			Expected: &Options{
				Version: true,
			},
		},
		{
			Args:         []string{"-debug"},
			TargetType:   checker.TargetTypeIsUnityProjectRootDirectory,
			RootDirAbs:   typedpath.NewRawPath("path", "to", "unity", "proj"),
			IgnoredPaths: []globs.Glob{},
			Expected: &Options{
				Version:        false,
				LogLevel:       logging.SeverityDebug,
				TargetType:     checker.TargetTypeIsUnityProjectRootDirectory,
				IgnoreDangling: false,
				IgnoreCase:     true,
				IgnoredPaths:   []globs.Glob{},
				RootDirAbs:     typedpath.NewRawPath("path", "to", "unity", "proj"),
			},
		},
		{
			Args:         []string{"-silent"},
			TargetType:   checker.TargetTypeIsUnityProjectRootDirectory,
			RootDirAbs:   typedpath.NewRawPath("path", "to", "unity", "proj"),
			IgnoredPaths: []globs.Glob{},
			Expected: &Options{
				Version:        false,
				LogLevel:       logging.SeverityWarn,
				TargetType:     checker.TargetTypeIsUnityProjectRootDirectory,
				IgnoreDangling: false,
				IgnoreCase:     true,
				IgnoredPaths:   []globs.Glob{},
				RootDirAbs:     typedpath.NewRawPath("path", "to", "unity", "proj"),
			},
		},
		{
			Args:         []string{"-silent", "-debug"},
			TargetType:   checker.TargetTypeIsUnityProjectRootDirectory,
			RootDirAbs:   typedpath.NewRawPath("path", "to", "unity", "proj"),
			IgnoredPaths: []globs.Glob{},
			Expected: &Options{
				Version:        false,
				LogLevel:       logging.SeverityDebug,
				TargetType:     checker.TargetTypeIsUnityProjectRootDirectory,
				IgnoreDangling: false,
				IgnoreCase:     true,
				IgnoredPaths:   []globs.Glob{},
				RootDirAbs:     typedpath.NewRawPath("path", "to", "unity", "proj"),
			},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c), func(t *testing.T) {
			stdin := strings.NewReader("")
			stdout := testutil.SpyWriteCloser(nil)
			stderr := testutil.SpyWriteCloser(nil)
			spyLogger := logging.SpyLogger()
			procInout := cli.ProcessInout{Stdin: stdin, Stdout: stdout, Stderr: stderr}

			buildOptions := NewBuilder(
				StubRootDirCompletion(c.RootDirAbs, nil),
				StubUnityProjectDetector(c.TargetType, nil),
				StubIgnoredPathBuilder(c.IgnoredPaths, nil),
				StubRootDirValidator(c.RootDirAbs, nil),
				spyLogger,
			)

			actual, err := buildOptions(c.Args, procInout)
			if err != nil {
				t.Log(spyLogger.Logs.String())
				t.Errorf("want nil, got %#v", err)
				return
			}

			if !reflect.DeepEqual(actual, c.Expected) {
				t.Log(spyLogger.Logs.String())
				t.Error(cmp.Diff(c.Expected, actual))
				return
			}
		})
	}
}
