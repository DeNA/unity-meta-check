package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/DeNA/unity-meta-check/report"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-junit/junit"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/DeNA/unity-meta-check/version"
	"io"
	"os"
	"time"
)

func main() {
	main := NewMain()
	exitStatus := main(os.Args[1:], cli.GetProcessInout(), cli.NewEnv())
	os.Exit(int(exitStatus))
}

func NewMain() cli.Command {
	return func(args []string, procInout cli.ProcessInout, env cli.Env) cli.ExitStatus {
		startTime := time.Now()

		opts, err := buildOptions(args, procInout)
		if err != nil {
			if err != flag.ErrHelp {
				_, _ = fmt.Fprintln(procInout.Stderr, err.Error())
			}
			return cli.ExitAbnormal
		}

		if opts.version {
			_, _ = fmt.Fprintln(procInout.Stdout, version.Version)
			return cli.ExitNormal
		}

		parse := report.NewParser()
		result := parse(io.TeeReader(procInout.Stdin, procInout.Stdout))

		if err := junit.WriteToFile(result, startTime, opts.outPath); err != nil {
			_, _ = fmt.Fprintln(procInout.Stderr, err.Error())
			return cli.ExitAbnormal
		}

		if !result.Empty() {
			return cli.ExitAbnormal
		}
		return cli.ExitNormal
	}
}

type options struct {
	version bool
	outPath typedpath.RawPath
}

func buildOptions(args []string, procInout cli.ProcessInout) (*options, error) {
	opts := &options{}

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
	flags.BoolVar(&opts.version, "version", false, "print version")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	if opts.version {
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
	opts.outPath = outPath
	return opts, nil
}
