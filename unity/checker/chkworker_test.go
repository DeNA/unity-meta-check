package checker

import (
	"github.com/DeNA/unity-meta-check/filecollector"
	"github.com/DeNA/unity-meta-check/unity"
	"github.com/DeNA/unity-meta-check/util/chanutil"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestNewCheckUnityProject(t *testing.T) {
	cases := map[string]struct {
		FoundPaths []filecollector.Entry
		IgnoreCase bool
		Expected   *CheckResult
	}{
		"no paths found (boundary)": {
			FoundPaths: []filecollector.Entry{},
			IgnoreCase: false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},
		"success case (representative)": {
			FoundPaths: []filecollector.Entry{
				{Path: "Assets/OK", IsDir: false},
				{Path: "Assets/OK.meta", IsDir: false},
			},
			IgnoreCase: false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},
		"several missing metas case (representative)": {
			FoundPaths: []filecollector.Entry{
				{Path: "Assets/NG", IsDir: false},
			},
			IgnoreCase: false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{"Assets/NG.meta"},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},
		"first missing metas w/ others case (representative)": {
			FoundPaths: []filecollector.Entry{
				{Path: "Assets/NG", IsDir: false},
				{Path: "Assets/OK", IsDir: false},
				{Path: "Assets/OK.meta", IsDir: false},
			},
			IgnoreCase: false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{"Assets/NG.meta"},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},
		"last missing metas w/ others case (representative)": {
			FoundPaths: []filecollector.Entry{
				{Path: "Assets/OK", IsDir: false},
				{Path: "Assets/OK.meta", IsDir: false},
				{Path: "Assets/NG", IsDir: false},
			},
			IgnoreCase: false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{"Assets/NG.meta"},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},
		"dangling meta (representative)": {
			FoundPaths: []filecollector.Entry{
				{Path: "Assets/NG.meta", IsDir: false},
			},
			IgnoreCase: false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{"Assets/NG.meta"},
			},
		},
		"last dangling meta w/ others (representative)": {
			FoundPaths: []filecollector.Entry{
				{Path: "Assets/OK", IsDir: false},
				{Path: "Assets/OK.meta", IsDir: false},
				{Path: "Assets/NG.meta", IsDir: false},
			},
			IgnoreCase: false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{"Assets/NG.meta"},
			},
		},
		"both missing dangling meta (representative)": {
			FoundPaths: []filecollector.Entry{
				{Path: "Assets/NG1", IsDir: false},
				{Path: "Assets/NG2.meta", IsDir: false},
			},
			IgnoreCase: false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{"Assets/NG1.meta"},
				DanglingMeta: []typedpath.SlashPath{"Assets/NG2.meta"},
			},
		},
		"directory contain only dangling meta (representative)": {
			FoundPaths: []filecollector.Entry{
				{Path: "Assets/Dangling", IsDir: true},
				{Path: "Assets/Dangling.meta", IsDir: false},
				{Path: "Assets/Dangling/Dangling1", IsDir: true},
				{Path: "Assets/Dangling/Dangling1.meta", IsDir: false},
				{Path: "Assets/Dangling/Dangling1/Dangling1.meta", IsDir: false},
				{Path: "Assets/Dangling/Dangling1/Dangling2.meta", IsDir: false},
				{Path: "Assets/Dangling/Dangling2", IsDir: true},
				{Path: "Assets/Dangling/Dangling2.meta", IsDir: false},
				{Path: "Assets/Dangling/Dangling2/Dangling3.meta", IsDir: false},
			},
			IgnoreCase: false,
			Expected: &CheckResult{
				MissingMeta: []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{
					"Assets/Dangling.meta",
					"Assets/Dangling/Dangling1.meta",
					"Assets/Dangling/Dangling1/Dangling1.meta",
					"Assets/Dangling/Dangling1/Dangling2.meta",
					"Assets/Dangling/Dangling2.meta",
					"Assets/Dangling/Dangling2/Dangling3.meta",
				},
			},
		},
		"needless meta (representative)": {
			FoundPaths: []filecollector.Entry{
				// NOTE: This path is not in Assets, so .meta is needless. But this case frequently happen by several
				//       reason e.g. the package including examples, but it is not used by the Unity project.
				//       So it should not be counted as a dangling.
				{Path: "foo/bar", IsDir: false},
				{Path: "foo/bar.meta", IsDir: false},
			},
			IgnoreCase: false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},

		"case-sensitive case (representative)": {
			FoundPaths: []filecollector.Entry{
				{Path: "Assets/NG1", IsDir: false},
				{Path: "Assets/ng1.meta", IsDir: false},
			},
			IgnoreCase: false,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{"Assets/NG1.meta"},
				DanglingMeta: []typedpath.SlashPath{"Assets/ng1.meta"},
			},
		},
		"case-insensitive case (representative)": {
			FoundPaths: []filecollector.Entry{
				{Path: "Assets/NG1", IsDir: false},
				{Path: "Assets/ng1.meta", IsDir: false},
			},
			IgnoreCase: true,
			Expected: &CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},
		// https://github.com/DeNA/unity-meta-check/issues/26
		"a directory contains a directory only contains normal file": {
			FoundPaths: []filecollector.Entry{
				{Path: "Assets/A", IsDir: true},
				{Path: "Assets/A.meta", IsDir: true},
				{Path: "Assets/A/B", IsDir: true},
				{Path: "Assets/A/B.meta", IsDir: false},
				{Path: "Assets/A/B/C.cs", IsDir: false},
				{Path: "Assets/A/B/C.cs.meta", IsDir: false},
				{Path: "Assets/A/D", IsDir: true},
				{Path: "Assets/A/D.meta", IsDir: false},
				{Path: "Assets/A/D/E", IsDir: true},
				{Path: "Assets/A/D/E.meta", IsDir: false},
			},
			IgnoreCase: false,
			Expected: &CheckResult{
				MissingMeta: []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{
					"Assets/A/D.meta",
					"Assets/A/D/E.meta",
				},
			},
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			spyLogger := logging.SpyLogger()
			entries := chanutil.FromSlice(c.FoundPaths)

			requiresMeta := unity.NewMetaNecessityInUnityProject([]typedpath.SlashPath{})
			check := NewCheckingWorker(requiresMeta, spyLogger)
			actual, err := check(".", c.IgnoreCase, entries)
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
