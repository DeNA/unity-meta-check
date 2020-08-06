package options

import (
	"github.com/DeNA/unity-meta-check/ignore"
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"os"
)

const IgnoreFileBasename typedpath.BaseName = ".meta-check-ignore"

type IgnoredPathsBuilder func(ignoreFilePath typedpath.RawPath, rootDirAbs typedpath.RawPath) ([]globs.Glob, error)

func NewIgnoredPathsBuilder(logger logging.Logger) IgnoredPathsBuilder {
	return func(ignoreFilePath typedpath.RawPath, rootDirAbs typedpath.RawPath) ([]globs.Glob, error) {
		ignoredPaths, err := ignore.ReadFile(getIgnoreFilePath(ignoreFilePath, rootDirAbs))
		if err != nil {
			// NOTE: If it is a default value, missing .meta-check-ignore is allowed because it is optional.
			//       Otherwise, treat as an error if specified ignoreFilePath is missing.
			if ignoreFilePath != "" || !os.IsNotExist(err) {
				return nil, err
			}
			logger.Info("no .meta-check-ignore, so ignored paths are empty")
			ignoredPaths = []globs.Glob{}
		}
		return ignoredPaths, nil
	}
}

func getIgnoreFilePath(path typedpath.RawPath, rootDirAbs typedpath.RawPath) typedpath.RawPath {
	if path == "" {
		return rootDirAbs.JoinBaseName(IgnoreFileBasename)
	} else {
		return path
	}
}
