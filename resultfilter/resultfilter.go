package resultfilter

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

type Options struct {
	IgnoreDangling bool
	IgnoredGlobs   []globs.Glob
	IgnoreCase     bool
}

type Filter func(result *checker.CheckResult, opts *Options) (*checker.CheckResult, error)

func NewFilter(logger logging.Logger) Filter {
	return func(result *checker.CheckResult, opts *Options) (*checker.CheckResult, error) {
		var ignored bool
		var matched globs.Glob
		var err error

		newMissingMeta := make([]typedpath.SlashPath, 0)
		for _, missingMeta := range result.MissingMeta {
			ignored, matched, err = globs.MatchAny(missingMeta, opts.IgnoredGlobs)
			if err != nil {
				return nil, err
			}
			if ignored {
				logger.Debug(fmt.Sprintf("ignored missing: %q by %q", missingMeta, matched))
				continue
			}
			newMissingMeta = append(newMissingMeta, missingMeta)
		}

		newDanglingMeta := make([]typedpath.SlashPath, 0)
		if !opts.IgnoreDangling {
			for _, danglingMeta := range result.DanglingMeta {
				ignored, matched, err = globs.MatchAny(danglingMeta, opts.IgnoredGlobs)
				if err != nil {
					return nil, err
				}
				if ignored {
					logger.Debug(fmt.Sprintf("ignored dangling: %q by %q", danglingMeta, matched))
					continue
				}
				newDanglingMeta = append(newDanglingMeta, danglingMeta)
			}
		}

		return checker.NewCheckResult(newMissingMeta, newDanglingMeta), nil
	}
}
