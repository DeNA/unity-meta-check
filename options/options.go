package options

import (
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

type Options struct {
	Version                   bool
	LogLevel                  logging.Severity
	TargetType                checker.TargetType
	IgnoreDangling            bool
	IgnoreCase                bool
	IgnoreSubmodulesAndNested bool
	IgnoredPaths              []globs.Glob
	RootDirAbs                typedpath.RawPath
}
