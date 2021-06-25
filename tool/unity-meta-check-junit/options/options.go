package options

import (
	"errors"
	"flag"
	"fmt"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

type Options struct {
	Version bool
	OutPath typedpath.RawPath
}

func BuildOptions(args []string, procInout cli.ProcessInout) (*Options, error) {
	opts := &Options{}

	flags := flag.NewFlagSet("unity-meta-check-junit", flag.ContinueOnError)
	flags.SetOutput(procInout.Stderr)
	flags.Usage = func() {
		_, _ = fmt.Fprint(procInout.Stderr, `usage: unity-meta-check-junit [<options>] [<path>]

Save a JUnit report file for the result from unity-meta-check via STDIN.

  <path>
        output path to write JUnit report

OPTIONS
`)
		flags.PrintDefaults()

		_, _ = fmt.Fprint(procInout.Stderr, `
EXAMPLE USAGES
  $ unity-meta-check <options> | unity-meta-check-junit path/to/junit-report.xml
  $ unity-meta-check <options> | unity-meta-check-junit path/to/junit-report.xml | <other-unity-meta-check-tool>
`)
	}
	flags.BoolVar(&opts.Version, "version", false, "print version")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	if opts.Version {
		return opts, nil
	}

	outPaths := flags.Args()
	if len(outPaths) < 1 {
		return nil, errors.New("must specify a file path to output JUnit report")
	}
	if len(outPaths) > 1 {
		return nil, errors.New("too much arguments")
	}

	outPath := typedpath.NewRawPathUnsafe(args[0])
	opts.OutPath = outPath
	return opts, nil
}
