package options

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/git"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/ostestable"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/pkg/errors"
	"path/filepath"
)

type RootDirCompletion func() (typedpath.RawPath, error)

func NewRootDirCompletion(gitRevParse git.RevParse, logger logging.Logger) RootDirCompletion {
	return func() (typedpath.RawPath, error) {
		assumedRootDir, err := gitRevParse(".", "--show-toplevel")
		if err != nil {
			return "", errors.Wrap(err, "rootDir not specified and seems not being in any git repositories")
		}

		logger.Debug(fmt.Sprintf("rootDir not specified, so assumed by $(git rev-parse --show-toplevel): %q", assumedRootDir))
		return typedpath.NewRawPathUnsafe(assumedRootDir), nil
	}
}

type RootDirAbsValidator func(unsafeRootDir typedpath.RawPath) (typedpath.RawPath, error)

func NewRootDirValidator(isDir ostestable.IsDir) RootDirAbsValidator {
	return func(unsafeRootDir typedpath.RawPath) (typedpath.RawPath, error) {
		rootDirAbsStr, err := filepath.Abs(string(unsafeRootDir))
		if err != nil {
			return "", err
		}

		rootDirAbs := typedpath.NewRawPathUnsafe(rootDirAbsStr)

		ok, err := isDir(rootDirAbs)
		if err != nil {
			return "", errors.Wrapf(err, "cannot check directory: %q", rootDirAbs)
		}
		if !ok {
			return "", fmt.Errorf("root directory must be a directory: %s", rootDirAbsStr)
		}

		return rootDirAbs, nil
	}
}
