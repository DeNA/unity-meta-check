package autofix

import (
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

func StubOptionsBuilderWithRootDirAbsAndRel(rootDirRel typedpath.RawPath) OptionsBuilder {
	return func(rootDirAbs typedpath.RawPath, allowedGlobs []globs.Glob) (*Options, error) {
		return &Options{
			RootDirAbs:   rootDirAbs,
			RootDirRel:   rootDirRel,
			AllowedGlobs: allowedGlobs,
		}, nil
	}
}
