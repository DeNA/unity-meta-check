package autofix

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/ostestable"
)

type AutoFixer func(result *checker.CheckResult, opts *Options) error

func NewAutoFixer(dryRun bool, getwd ostestable.Getwd, detectMetaType MetaTypeDetector, createMeta MetaCreator, removeMeta MetaRemover, logger logging.Logger) AutoFixer {
	return func(result *checker.CheckResult, opts *Options) error {
		if result.Empty() {
			logger.Info("no missing or dangling .meta. nothing to do")
			return nil
		}

		rawCwd, err := getwd()
		if err != nil {
			return err
		}
		cwd := rawCwd.ToSlash()

		for _, missingMeta := range result.MissingMeta {
			ok, matched, err := globs.MatchAny(missingMeta, opts.AllowedGlobs, cwd)
			if err != nil {
				return err
			}
			if !ok {
				logger.Info(fmt.Sprintf("autofix skipped because: no specified patterns match to the missing .meta: %s", missingMeta))
				continue
			}
			logger.Debug(fmt.Sprintf("try to autofix missing .meta: %s (matched by %q)", missingMeta, matched))

			missingMetaAbs := opts.RootDirAbs.JoinRawPath(missingMeta.ToRaw())
			metaType, err := detectMetaType(missingMetaAbs)
			if err != nil {
				logger.Warn(fmt.Sprintf("generating .meta skipped because: %s", err.Error()))
				continue
			}
			logger.Debug(fmt.Sprintf("meta type detected: %q for %s (matched by %q)", metaType, missingMeta, matched))

			if err := createMeta(metaType, missingMetaAbs); err != nil {
				return err
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
				return err
			}
			if !ok {
				logger.Info(fmt.Sprintf("autofix skipped because: no specified patterns match to the dangling .meta: %s", danglingMeta))
				continue
			}
			logger.Debug(fmt.Sprintf("try to autofix dangling .meta: %s (matched by %q)", danglingMeta, matched))

			danglingMetaAbs := opts.RootDirAbs.JoinRawPath(danglingMeta.ToRaw())
			if err := removeMeta(danglingMetaAbs); err != nil {
				return err
			}

			danglingMetaRel := opts.RootDirRel.JoinRawPath(danglingMeta.ToRaw())
			if dryRun {
				logger.Info(fmt.Sprintf("would remove: %s (matched by %q)", danglingMetaRel, matched))
			} else {
				logger.Info(fmt.Sprintf("removed: %s (matched by %q)", danglingMetaRel, matched))
			}
		}

		return nil
	}
}
