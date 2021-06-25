package autofix

import (
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/ostestable"
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

type Options struct {
	RootDirAbs   typedpath.RawPath
	RootDirRel   typedpath.RawPath
	AllowedGlobs []globs.Glob
}

type OptionsBuilder func(rootDirAbs typedpath.RawPath, allowedGlobs []globs.Glob) (*Options, error)

func NewOptionsBuilder(getwd ostestable.Getwd) OptionsBuilder {
	return func(rootDirAbs typedpath.RawPath, allowedGlobs []globs.Glob) (*Options, error) {
		cwdAbs, err := getwd()
		if err != nil {
			return nil, err
		}

		rootDirRel, err := cwdAbs.Rel(rootDirAbs)
		if err != nil {
			return nil, err
		}

		return &Options{
			RootDirAbs:   rootDirAbs,
			RootDirRel:   rootDirRel,
			AllowedGlobs: allowedGlobs,
		}, nil
	}
}
