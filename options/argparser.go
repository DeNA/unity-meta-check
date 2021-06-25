package options

import (
	"flag"
	"fmt"
	"github.com/DeNA/unity-meta-check/git"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/ostestable"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/pkg/errors"
)

type ArgParser func(args []string, procInout cli.ProcessInout) (*Options, error)

var _ ArgParser = ParseArgs

func ParseArgs(args []string, procInout cli.ProcessInout) (*Options, error) {
	adhocLogger := logging.NewLogger(logging.SeverityWarn, procInout.Stderr)
	builder := NewBuilder(
		NewRootDirCompletion(git.NewRevParse(adhocLogger), adhocLogger),
		NewUnityProjectDetector(adhocLogger),
		NewIgnoredGlobsBuilder(adhocLogger),
		NewRootDirValidator(ostestable.NewIsDir()),
		adhocLogger,
	)
	opts, err := builder(args, procInout)
	if err != nil {
		return nil, err
	}
	return opts, nil
}

func NewBuilder(
	detectRootDir RootDirCompletion,
	detectUnityProj UnityProjectDetector,
	buildIgnoredGlobs IgnoredGlobsBuilder,
	validateRootDirAbs RootDirAbsValidator,
	logger logging.Logger,
) ArgParser {
	return func(args []string, procInout cli.ProcessInout) (*Options, error) {
		var opts Options
		flags := flag.NewFlagSet("unity-meta-check", flag.ContinueOnError)
		flags.Usage = func() {
			_, _ = fmt.Fprint(procInout.Stderr, `usage: unity-meta-check [<options>] [<path>]

Check missing or dangling .meta files.

  <path>
        root directory of your Unity project or UPM package to check (default "$(git rev-parse --show-toplevel)")

OPTIONS
`)
			flags.PrintDefaults()
			_, _ = fmt.Fprint(procInout.Stderr, `
EXAMPLE USAGES
  $ cd path/to/UnityProject
  $ unity-meta-check -silent

  $ cd path/to/any/dir
  $ unity-meta-check -silent -upm-package path/to/MyUPMPackage
  $ unity-meta-check -silent -unity-project-sub-dir path/to/UnityProject/Assets/Sub/Dir

EXAMPLE USAGES WITH OTHER TOOLS
  $ unity-meta-check <options> | unity-meta-check-junit path/to/junit.xml
  $ unity-meta-check <options> | unity-meta-check-github-pr-comment <options>
  $ unity-meta-check <options> | unity-meta-check-junit path/to/junit.xml | unity-meta-check-github-pr-comment <options>
`)
		}
		var noIgnoreCase, unityProj, upmPkg, unityProjSubDir, silent, debug, ignoreSubmodulesAndNested bool
		var ignoreFilePath string
		flags.BoolVar(&opts.Version, "version", false, "print version")
		flags.StringVar(&ignoreFilePath, "ignore-file", "", "path to .meta-check-ignore")
		flags.BoolVar(&upmPkg, "upm-package", false, "check as UPM package (same meaning of -unity-project-sub-dir)")
		flags.BoolVar(&unityProj, "unity-project", false, "check as Unity project")
		flags.BoolVar(&unityProjSubDir, "unity-project-sub-dir", false, "check as sub directory of Unity project")
		flags.BoolVar(&silent, "silent", false, "set log level to WARN (default INFO)")
		flags.BoolVar(&debug, "debug", false, "set log level to DEBUG (default INFO)")
		flags.BoolVar(&noIgnoreCase, "no-ignore-case", false, "treat case of file paths")
		flags.BoolVar(&opts.IgnoreDangling, "ignore-dangling", false, "ignore dangling .meta")
		flags.BoolVar(&ignoreSubmodulesAndNested, "ignore-submodules", false, "ignore git submodules and nesting repositories (this is RECOMMENDED but not enabled by default because it can cause to miss problems in submodules or nesting repositories)")
		flags.SetOutput(procInout.Stderr)

		if err := flags.Parse(args); err != nil {
			return nil, err
		}

		if opts.Version {
			return &opts, nil
		}

		logLevel := cli.GetLogLevel(debug, silent)
		opts.LogLevel = logLevel

		var unsafeRootDir typedpath.RawPath
		switch flags.NArg() {
		case 0:
			var err error
			unsafeRootDir, err = detectRootDir()
			if err != nil {
				return nil, err
			}
		case 1:
			unsafeRootDir = typedpath.NewRawPathUnsafe(flags.Arg(0))
		default:
			return nil, fmt.Errorf("want exactly 1 argument as a directory path to check, got %d arguments: %#v", flags.NArg(), flags.Args())
		}

		rootDirAbs, err := validateRootDirAbs(unsafeRootDir)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid root directory: %q", unsafeRootDir)
		}

		opts.RootDirAbs = rootDirAbs

		opts.IgnoreCase = !noIgnoreCase

		explicitTargetType, explicitTargetTypeOk, err := validateTargetTypeFlags(unityProj, unityProjSubDir, upmPkg)
		if err != nil {
			return nil, err
		}

		var targetType checker.TargetType
		if explicitTargetTypeOk {
			targetType = explicitTargetType
		} else {
			logger.Info("none of -upm-package and -unity-project and -unity-project-sub-dir was specified, so try to detect it.")

			targetType, err = detectUnityProj(rootDirAbs)
			if err != nil {
				return nil, err
			}
		}
		opts.TargetType = targetType

		ignoredPaths, err := buildIgnoredGlobs(typedpath.NewRawPathUnsafe(ignoreFilePath), rootDirAbs)
		if err != nil {
			return nil, err
		}
		opts.IgnoredPaths = ignoredPaths

		opts.IgnoreSubmodulesAndNested = ignoreSubmodulesAndNested

		return &opts, nil
	}
}

func validateTargetTypeFlags(unityProj, unityProjSubDir, upmPkg bool) (checker.TargetType, bool, error) {
	notUnityProj := upmPkg || unityProjSubDir

	if notUnityProj && unityProj {
		return "", false, errors.New("must specify one of -upm-package or -unity-project or -unity-project-sub-dir")
	}

	if notUnityProj {
		return checker.TargetTypeIsUnityProjectSubDirectory, true, nil
	}

	if unityProj {
		return checker.TargetTypeIsUnityProjectRootDirectory, true, nil
	}

	return "", false, nil
}
