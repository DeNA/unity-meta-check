package autofix

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/ostestable"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"path/filepath"
	"testing"
)

func TestNewAutoFixer(t *testing.T) {
	cases := []struct {
		DryRun       bool
		RootDirRel   typedpath.RawPath
		Target       *checker.CheckResult
		AllowedGlobs []globs.Glob
		ExpectedErr  bool
	}{
		{
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
		},
		{
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
		},
		{
			DryRun:     true,
			RootDirRel: typedpath.NewRawPath("testdata", "InvalidProject"),
			Target: checker.NewCheckResult(
				[]typedpath.SlashPath{"Assets/Missing.meta"},
				[]typedpath.SlashPath{"Assets/Dangling.meta"},
			),
			AllowedGlobs: []globs.Glob{
				globs.Glob(filepath.Join("Assets", "*")),
			},
			ExpectedErr: false,
		},
		{
			DryRun:     false,
			RootDirRel: typedpath.NewRawPath("testdata", "InvalidProject"),
			Target: checker.NewCheckResult(
				[]typedpath.SlashPath{"Assets/Missing.meta"},
				[]typedpath.SlashPath{"Assets/Dangling.meta"},
			),
			AllowedGlobs: []globs.Glob{
				globs.Glob(filepath.Join("Assets", "*")),
			},
			ExpectedErr: false,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%q, %v, %v", c.RootDirRel, c.Target, c.AllowedGlobs), func(t *testing.T) {
			metaTypeDetector := StubMetaTypeDetector(MetaTypeTextScriptImporter, nil)
			metaCreator := StubMetaCreator(nil)
			metaRemover := StubMetaRemover(nil)
			spyLogger := logging.SpyLogger()

			autofix := NewAutoFixer(c.DryRun, ostestable.NewGetwd(), metaTypeDetector, metaCreator, metaRemover, spyLogger)

			err := autofix(c.Target, &Options{
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
					return
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
