package options

import (
	"errors"
	"flag"
	"fmt"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/globs"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	Version      bool
	LogLevel     logging.Severity
	DryRun       bool
	FixMissing   bool
	FixDangling  bool
	AllowedGlobs []globs.Glob
	RootDirAbs   typedpath.RawPath
}

func Build(args []string, procInout cli.ProcessInout) (*Options, error) {
	opts := Options{}

	flags := flag.NewFlagSet("unity-meta-autofix", flag.ContinueOnError)
	flags.SetOutput(procInout.Stderr)
	flags.Usage = func() {
		_, _ = fmt.Fprint(procInout.Stderr, `usage: unity-meta-autofix [<Options>] <pattern> [<pattern>...]

Fix missing or dangling .meta. Currently autofix is only limited support.

ARGUMENTS
  <pattern>
        glob pattern to path where autofix allowed on

OPTIONS
`)
		flags.PrintDefaults()

		_, _ = fmt.Fprint(procInout.Stderr, `
EXAMPLE USAGES
  $ unity-meta-check <Options> | unity-meta-autofix -dry-run -fix-missing -fix-dangling path/to/autofix
  $ unity-meta-check <Options> | unity-meta-autofix <Options> | <other-unity-meta-check-tool>
`)
	}

	var silent, debug bool
	var rootDir string
	flags.BoolVar(&opts.Version, "version", false, "print Version")
	flags.BoolVar(&debug, "debug", false, "set log level to DEBUG (default INFO)")
	flags.BoolVar(&silent, "silent", false, "set log level to WARN (default INFO)")
	flags.BoolVar(&opts.DryRun, "dry-run", false, "dry run")
	flags.BoolVar(&opts.FixMissing, "fix-missing", false, "fix missing .meta")
	flags.BoolVar(&opts.FixDangling, "fix-dangling", false, "fix dangling .meta")
	flags.StringVar(&rootDir, "root-dir", ".", "directory path to where unity-meta-check checked at")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	if opts.Version {
		return &opts, nil
	}

	if !opts.FixMissing && !opts.FixDangling {
		return nil, errors.New("must specify either or both -fix-missing and -fix-dangling")
	}

	targetPaths := flags.Args()
	if len(targetPaths) == 0 {
		return nil, errors.New("must specify at least one target path")
	}
	opts.AllowedGlobs = make([]globs.Glob, len(targetPaths))
	for i, targetPath := range targetPaths {
		opts.AllowedGlobs[i] = globs.Glob(strings.Trim(filepath.ToSlash(targetPath), "/"))
	}

	rootDirAbs, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, err
	}
	stat, err := os.Stat(rootDirAbs)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("rootDir must be a directory: %s", rootDirAbs)
	}
	opts.RootDirAbs = typedpath.NewRawPathUnsafe(rootDirAbs)

	opts.LogLevel = cli.GetLogLevel(debug, silent)

	return &opts, nil
}
