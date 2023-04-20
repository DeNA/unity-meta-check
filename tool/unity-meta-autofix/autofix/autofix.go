package autofix

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/ostestable"
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

type AutoFixer func(result *checker.CheckResult, opts *Options) (*checker.CheckResult, error)

func NewAutoFixer(dryRun bool, getwd ostestable.Getwd, detectMetaType MetaTypeDetector, createMeta MetaCreator, removeMeta MetaRemover, logger logging.Logger) AutoFixer {
	return func(result *checker.CheckResult, opts *Options) (*checker.CheckResult, error) {
		if result.Empty() {
			logger.Info("no missing or dangling .meta. nothing to do")
			return result, nil
		}

		rawCwd, err := getwd()
		if err != nil {
			return nil, err
		}
		cwd := rawCwd.ToSlash()

		skippedMissing := make([]typedpath.SlashPath, 0)
		skippedDangling := make([]typedpath.SlashPath, 0)
		for _, missingMeta := range result.MissingMeta {
			ok, matched, err := globs.MatchAny(missingMeta, opts.AllowedGlobs, cwd)
			if err != nil {
				return nil, err
			}
			if !ok {
				logger.Info(fmt.Sprintf("autofix skipped because: no specified patterns match to the missing .meta: %s", missingMeta))
				skippedMissing = append(skippedMissing, missingMeta)
				continue
			}
			logger.Debug(fmt.Sprintf("try to autofix missing .meta: %s (matched by %q)", missingMeta, matched))

			missingMetaAbs := opts.RootDirAbs.JoinRawPath(missingMeta.ToRaw())
			metaType, err := detectMetaType(missingMetaAbs)
			if err != nil {
				logger.Warn(fmt.Sprintf("generating .meta skipped because: %s", err.Error()))
				skippedMissing = append(skippedMissing, missingMeta)
				continue
			}
			logger.Debug(fmt.Sprintf("meta type detected: %q for %s (matched by %q)", metaType, missingMeta, matched))

			if err := createMeta(metaType, missingMetaAbs); err != nil {
				return nil, err
			}

			missingMetaRel := opts.RootDirRel.JoinRawPath(missingMeta.ToRaw())
			if dryRun {
				logger.Info(fmt.Sprintf("would generate: %s (matched by %q)", missingMetaRel, matched))
			} else {
				logger.Info(fmt.Sprintf("generated: %s (matched by %q)", missingMetaRel, matched))
			}
		}

		for _, danglingMeta := range result.DanglingMeta {
			ok, matched, err := globs.MatchAny(danglingMeta, opts.AllowedGlobs, cwd)
			if err != nil {
				return nil, err
			}
			if !ok {
				logger.Info(fmt.Sprintf("autofix skipped because: no specified patterns match to the dangling .meta: %s", danglingMeta))
				skippedDangling = append(skippedDangling, danglingMeta)
				continue
			}
			logger.Debug(fmt.Sprintf("try to autofix dangling .meta: %s (matched by %q)", danglingMeta, matched))

			danglingMetaAbs := opts.RootDirAbs.JoinRawPath(danglingMeta.ToRaw())
			if err := removeMeta(danglingMetaAbs); err != nil {
				return nil, err
			}

			danglingMetaRel := opts.RootDirRel.JoinRawPath(danglingMeta.ToRaw())
			if dryRun {
				logger.Info(fmt.Sprintf("would remove: %s (matched by %q)", danglingMetaRel, matched))
			} else {
				logger.Info(fmt.Sprintf("removed: %s (matched by %q)", danglingMetaRel, matched))
			}
		}

		return &checker.CheckResult{
			MissingMeta:  skippedMissing,
			DanglingMeta: skippedDangling,
		}, nil
	}
}
