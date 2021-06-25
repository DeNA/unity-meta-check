package options

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"strings"
)

func StubRootDirCompletion(path typedpath.RawPath, err error) RootDirCompletion {
	return func() (typedpath.RawPath, error) {
		return path, err
	}
}

func StubRootDirValidator(path typedpath.RawPath, err error) RootDirAbsValidator {
	return func(_ typedpath.RawPath) (typedpath.RawPath, error) {
		return path, err
	}
}

func FakeRootDirValidator(cwd typedpath.RawPath) RootDirAbsValidator {
	return func(unsafeRootDir typedpath.RawPath) (typedpath.RawPath, error) {
		if strings.HasPrefix(string(unsafeRootDir), "/") {
			return unsafeRootDir, nil
		}

		return cwd.JoinRawPath(unsafeRootDir), nil
	}
}
