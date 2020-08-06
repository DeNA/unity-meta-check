package options

import (
	"flag"
	"fmt"
	"github.com/DeNA/unity-meta-check/git"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

type ArgParser func(args []string, procInout cli.ProcessInout) (*Options, error)

var _ ArgParser = ParseArgs

func ParseArgs(args []string, procInout cli.ProcessInout) (*Options, error) {
	adhocLogger := logging.NewLogger(logging.SeverityWarn, procInout.Stderr)
	builder := NewBuilder(
		NewRootDirDetector(git.NewRevParse(adhocLogger), adhocLogger),
		NewUnityProjectDetector(adhocLogger),
		NewIgnoredPathsBuilder(adhocLogger),
	)
	opts, err := builder(args, procInout)
	if err != nil {
		return nil, err
	}
	return opts, nil
}

func NewBuilder(detectRootDir RootDirDetector, detectUnityProj UnityProjectDetector, buildIgnoredPaths IgnoredPathsBuilder) ArgParser {
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

		rootDirAbs, err := detectRootDir(flags.Args())
		if err != nil {
			return nil, err
		}
		opts.RootDirAbs = rootDirAbs

		opts.IgnoreCase = !noIgnoreCase

		isUnityProj, err := detectUnityProj(unityProj, upmPkg, unityProjSubDir, rootDirAbs)
		if err != nil {
			return nil, err
		}
		if isUnityProj {
			opts.TargetType = checker.TargetTypeIsUnityProjectRootDirectory
		} else {
			opts.TargetType = checker.TargetTypeIsUnityProjectSubDirectory
		}

		ignoredPaths, err := buildIgnoredPaths(typedpath.NewRawPathUnsafe(ignoreFilePath), rootDirAbs)
		if err != nil {
			return nil, err
		}
		opts.IgnoredPaths = ignoredPaths

		opts.IgnoreSubmodulesAndNested = ignoreSubmodulesAndNested

		return &opts, nil
	}
}
