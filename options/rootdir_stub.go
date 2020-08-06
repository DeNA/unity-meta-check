package options

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

func StubRootDirDetector(path typedpath.RawPath, err error) RootDirDetector {
	return func(_ []string) (typedpath.RawPath, error) {
		return path, err
	}
}
