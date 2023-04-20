package autofix

import (
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/ostestable"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewAutoFixer(t *testing.T) {
	cases := map[string]struct {
		DryRun          bool
		RootDirRel      typedpath.RawPath
		Target          *checker.CheckResult
		AllowedGlobs    []globs.Glob
		ExpectedErr     bool
		ExpectedSkipped *checker.CheckResult
	}{
		"dry-run & no missing or dangling metas": {
			DryRun:     true,
			RootDirRel: typedpath.NewRawPath("testdata", "ValidProject"),
			Target: checker.NewCheckResult(
				[]typedpath.SlashPath{},
				[]typedpath.SlashPath{},
			),
			AllowedGlobs: []globs.Glob{
				globs.Glob(filepath.Join("Assets", "*")),
			},
			ExpectedErr: false,
			ExpectedSkipped: &checker.CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},
		"not dry-run & no missing or dangling metas": {
			DryRun:     false,
			RootDirRel: typedpath.NewRawPath("testdata", "ValidProject"),
			Target: checker.NewCheckResult(
				[]typedpath.SlashPath{},
				[]typedpath.SlashPath{},
			),
			AllowedGlobs: []globs.Glob{
				globs.Glob(filepath.Join("Assets", "*")),
			},
			ExpectedErr: false,
			ExpectedSkipped: &checker.CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},
		"dry-run & several missing or dangling metas": {
			DryRun:     true,
			RootDirRel: typedpath.NewRawPath("testdata", "InvalidProject"),
			Target: checker.NewCheckResult(
				[]typedpath.SlashPath{"Assets/Missing.meta"},
				[]typedpath.SlashPath{"Assets/Dangling.meta"},
			),
			AllowedGlobs: []globs.Glob{
				globs.Glob("Assets/*"),
			},
			ExpectedErr: false,
			ExpectedSkipped: &checker.CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},
		"not dry-run & several missing or dangling metas": {
			DryRun:     false,
			RootDirRel: typedpath.NewRawPath("testdata", "InvalidProject"),
			Target: checker.NewCheckResult(
				[]typedpath.SlashPath{"Assets/Missing.meta"},
				[]typedpath.SlashPath{"Assets/Dangling.meta"},
			),
			AllowedGlobs: []globs.Glob{
				globs.Glob("Assets/*"),
			},
			ExpectedErr: false,
			ExpectedSkipped: &checker.CheckResult{
				MissingMeta:  []typedpath.SlashPath{},
				DanglingMeta: []typedpath.SlashPath{},
			},
		},
		"several skipped": {
			DryRun:     false,
			RootDirRel: typedpath.NewRawPath("testdata", "InvalidProject"),
			Target: checker.NewCheckResult(
				[]typedpath.SlashPath{"Assets/Missing.meta"},
				[]typedpath.SlashPath{"Assets/Dangling.meta"},
			),
			AllowedGlobs: []globs.Glob{},
			ExpectedErr:  false,
			ExpectedSkipped: &checker.CheckResult{
				MissingMeta:  []typedpath.SlashPath{"Assets/Missing.meta"},
				DanglingMeta: []typedpath.SlashPath{"Assets/Dangling.meta"},
			},
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			metaTypeDetector := StubMetaTypeDetector(MetaTypeTextScriptImporter, nil)
			metaCreator := StubMetaCreator(nil)
			metaRemover := StubMetaRemover(nil)
			spyLogger := logging.SpyLogger()

			autofix := NewAutoFixer(c.DryRun, ostestable.NewGetwd(), metaTypeDetector, metaCreator, metaRemover, spyLogger)

			skipped, err := autofix(c.Target, &Options{
				RootDirAbs:   rootDirAbs(c.RootDirRel),
				RootDirRel:   c.RootDirRel,
				AllowedGlobs: c.AllowedGlobs,
			})

			if c.ExpectedErr {
				if err == nil {
					t.Log(spyLogger.Logs.String())
					t.Errorf("want error, got nil")
					return
				}
			} else {
				if err != nil {
					t.Log(spyLogger.Logs.String())
					t.Errorf("want nil, got %#v", err)
				}

				if !reflect.DeepEqual(skipped, c.ExpectedSkipped) {
					t.Log(spyLogger.Logs.String())
					t.Error(cmp.Diff(c.ExpectedSkipped, skipped))
				}
			}
		})
	}
}

func rootDirAbs(rootDirRel typedpath.RawPath) typedpath.RawPath {
	cwdAbs, err := typedpath.Getwd()
	if err != nil {
		panic(err.Error())
	}
	return cwdAbs.JoinRawPath(rootDirRel)
}
