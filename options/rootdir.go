package options

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/git"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"os"
	"path/filepath"
)

type RootDirDetector func(nonFlaggedArgs []string) (typedpath.RawPath, error)

func NewRootDirDetector(gitRevParse git.RevParse, logger logging.Logger) RootDirDetector {
	return func(nonFlaggedArgs []string) (typedpath.RawPath, error) {
		var rootDir string
		if len(nonFlaggedArgs) == 0 {
			assumedRootDir, err := gitRevParse(".", "--show-toplevel")
			if err == nil {
				rootDir = assumedRootDir
				logger.Info(fmt.Sprintf("rootDir not specified, so assumed by $(git rev-parse --show-toplevel): %s", rootDir))
			} else {
				// NOTE: This is an fallback
				rootDir = "."
				logger.Info(fmt.Sprintf("rootDir not specified and seems not in any git repositories, so using the fallback value: %s", rootDir))
			}
		} else if len(nonFlaggedArgs) == 1 {
			rootDir = nonFlaggedArgs[0]
			logger.Info(fmt.Sprintf("rootDir specified by arguments: %s", rootDir))
		} else {
			return "", fmt.Errorf("must be only 1 path, but given %d paths", len(nonFlaggedArgs))
		}

		rootDirAbs, err := filepath.Abs(rootDir)
		if err != nil {
			return "", err
		}

		stat, err := os.Stat(rootDirAbs)
		if err != nil {
			return "", err
		}
		if !stat.IsDir() {
			return "", fmt.Errorf("not a directory: %q", rootDirAbs)
		}

		return typedpath.RawPath(rootDirAbs), err
	}
}
