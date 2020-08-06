package checker

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/unity"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/pathchan"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestNewCheckUnityProject(t *testing.T) {
	cases := []struct {
		FoundPaths    []typedpath.SlashPath
		IgnoreCase    bool
		Expected      *CheckResult
	}{
		{
			FoundPaths: []typedpath.SlashPath{},
			IgnoreCase: false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},
		{
			FoundPaths: []typedpath.SlashPath{
				"Assets/OK",
				"Assets/OK.meta",
			},
			IgnoreCase: false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},
		{
			FoundPaths: []typedpath.SlashPath{
				"Assets/NG",
			},
			IgnoreCase:    false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{"Assets/NG.meta"},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},
		{
			FoundPaths: []typedpath.SlashPath{
				"Assets/NG",
				"Assets/OK",
				"Assets/OK.meta",
			},
			IgnoreCase:    false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{"Assets/NG.meta"},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},
		{
			FoundPaths: []typedpath.SlashPath{
				"Assets/OK",
				"Assets/OK.meta",
				"Assets/NG",
			},
			IgnoreCase:    false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{"Assets/NG.meta"},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},

		{
			FoundPaths: []typedpath.SlashPath{
				"Assets/NG.meta",
			},
			IgnoreCase:    false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{"Assets/NG.meta"},
			},
		},
		{
			FoundPaths: []typedpath.SlashPath{
				"Assets/OK",
				"Assets/OK.meta",
				"Assets/NG.meta",
			},
			IgnoreCase:    false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{"Assets/NG.meta"},
			},
		},
		{
			FoundPaths: []typedpath.SlashPath{
				"Assets/NG1",
				"Assets/NG2.meta",
			},
			IgnoreCase:    false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{"Assets/NG1.meta"},
				DanglingMeta: []typedpath.SlashPath{"Assets/NG2.meta"},
			},
		},
		{
			FoundPaths: []typedpath.SlashPath{
				// NOTE: This path is not in Assets, so .meta is needless. But this case frequently happen by several
				//       reason e.g. the package including examples, but it is not used by the Unity project.
				//       So it should not be counted as a dangling.
				"foo/bar",
				"foo/bar.meta",
			},
			IgnoreCase:    false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},

		// Ignore case true.
		{
			FoundPaths: []typedpath.SlashPath{
				"Assets/NG1",
				"Assets/ng1.meta",
			},
			IgnoreCase:    true,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v IgnoreCase=%v", c.FoundPaths, c.IgnoreCase), func(t *testing.T) {
			spyLogger := logging.SpyLogger()
			foundPaths := pathchan.FromSlice(c.FoundPaths)

			requiresMeta := unity.NewMetaNecessityInUnityProject([]typedpath.SlashPath{})
			check := NewCheckingWorker(requiresMeta, spyLogger)
			actual, err := check(c.IgnoreCase, foundPaths)
			if err != nil {
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
