package unity

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestNewFindPackages(t *testing.T) {
	cwd, err := typedpath.Getwd()
	if err != nil {
		panic(err.Error())
	}

	cases := []struct{
		RootDirAbs typedpath.RawPath
		Expected   []*FoundPackage
	} {
		{
			cwd.JoinRawPath(typedpath.NewRawPath("testdata", "EmptyManifest")),
			[]*FoundPackage{},
		},
		{
			cwd.JoinRawPath(typedpath.NewRawPath("testdata", "NoLocalsManifest")),
			[]*FoundPackage{
				{
					FilePrefix: false,
					AbsPath:    cwd.JoinRawPath(typedpath.NewRawPath("testdata", "NoLocalsManifest", "Packages", "com.example.exists")),
					RelPath:    typedpath.NewRawPath("Packages", "com.example.exists"),
				},
			},
		},
		{
			cwd.JoinRawPath(typedpath.NewRawPath("testdata", "LocalsManifest")),
			[]*FoundPackage{
				{
					FilePrefix: true,
					AbsPath:    cwd.JoinRawPath(typedpath.NewRawPath("testdata", "LocalsManifest", "LocalPackages", "com.example.local")),
					RelPath:    typedpath.NewRawPath("LocalPackages", "com.example.local"),
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s/Packages/manifest.json -> %v", c.RootDirAbs, c.Expected), func(t *testing.T) {
			spyLogger := logging.SpyLogger()
			findPackages := NewFindPackages(spyLogger)

			actual, err := findPackages(c.RootDirAbs)
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
