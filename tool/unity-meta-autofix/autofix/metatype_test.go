package autofix

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/util/ostestable"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"testing"
)

func TestNewMetaTypeDetectorValid(t *testing.T) {
	cases := []struct {
		MissingMeta typedpath.RawPath
		IsDir       bool
		Expected    MetaType
	}{
		{
			MissingMeta: typedpath.NewRawPath("Assets", "Dir.meta"),
			IsDir:       true,
			Expected:    MetaTypeDefaultImporterFolder,
		},
		{
			MissingMeta: typedpath.NewRawPath("Assets", "file.json.meta"),
			IsDir:       false,
			Expected: MetaTypeTextScriptImporter,
		},
		{
			MissingMeta: typedpath.NewRawPath("Assets", "FILE.JSON.meta"),
			IsDir:       false,
			Expected: MetaTypeTextScriptImporter,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%q (isDir=%t) -> %q", c.MissingMeta, c.IsDir, c.Expected), func(t *testing.T) {
			detectMetaType := NewMetaTypeDetector(ostestable.StubIsDir(c.IsDir, nil))

			actual, err := detectMetaType(c.MissingMeta)
			if err != nil {
				t.Errorf("want nil, got %#v", err)
				return
			}

			if actual != c.Expected {
				t.Errorf("want %q, got %q", c.Expected, actual)
				return
			}
		})
	}
}

func TestNewMetaTypeDetectorInvalid(t *testing.T) {
	cases := []struct {
		MissingMeta typedpath.RawPath
		IsDir       bool
	}{
		{
			MissingMeta: typedpath.NewRawPath("Assets", "Some.pkg.meta"),
			IsDir:       true,
		},
		{
			MissingMeta: typedpath.NewRawPath("Assets", "file.unknown.meta"),
			IsDir:       false,
		},
		{
			MissingMeta: typedpath.NewRawPath("Assets", "no-extension.meta"),
			IsDir:       false,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%q (isDir=%t) -> error", c.MissingMeta, c.IsDir), func(t *testing.T) {
			detectMetaType := NewMetaTypeDetector(ostestable.StubIsDir(c.IsDir, nil))

			_, err := detectMetaType(c.MissingMeta)
			if err == nil {
				t.Errorf("want error, got nil")
				return
			}
		})
	}
}
