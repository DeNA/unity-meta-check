package options

import (
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

func StubIgnoredPathBuilder(result []globs.Glob, err error) IgnoredPathsBuilder {
	return func(typedpath.RawPath, typedpath.RawPath) ([]globs.Glob, error) {
		return result, err
	}
}
