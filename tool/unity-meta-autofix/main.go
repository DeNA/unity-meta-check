package main

import (
	"flag"
	"fmt"
	"github.com/DeNA/unity-meta-check/report"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-autofix/autofix"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-autofix/options"
	"github.com/DeNA/unity-meta-check/unity/meta"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/ostestable"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"github.com/DeNA/unity-meta-check/version"
	"os"
)

func main() {
	main := NewMain()
	exitStatus := main(os.Args[1:], cli.GetProcessInout(), cli.NewEnv())
	os.Exit(int(exitStatus))
}

func NewMain() cli.Command {
	return func(args []string, procInout cli.ProcessInout, env cli.Env) cli.ExitStatus {
		opts, err := options.Build(args, procInout)
		if err != nil {
			if err != flag.ErrHelp {
				_, _ = fmt.Fprintln(procInout.Stderr, err.Error())
			}
			return cli.ExitAbnormal
		}

		if opts.Version {
			_, _ = fmt.Fprintln(procInout.Stdout, version.Version)
			return cli.ExitNormal
		}

		parse := report.NewParser()
		result := parse(procInout.Stdin)

		logger := logging.NewLogger(opts.LogLevel, procInout.Stderr)
		cwdAbs, err := typedpath.Getwd()
		if err != nil {
			_, _ = fmt.Fprintln(procInout.Stderr, err.Error())
			return cli.ExitAbnormal
		}
		rootDirRel, err := cwdAbs.Rel(opts.RootDirAbs)
		if err != nil {
			_, _ = fmt.Fprintln(procInout.Stderr, err.Error())
			return cli.ExitAbnormal
		}

		autofixFunc := autofix.NewAutoFixer(
			opts.DryRun,
			autofix.NewMetaTypeDetector(ostestable.NewIsDir()),
			autofix.NewMetaCreator(opts.DryRun, meta.RandomGUIDGenerator(), logger),
			autofix.NewMetaRemover(opts.DryRun),
			logger,
		)
		autofixOpts := autofix.NewOptions(opts.RootDirAbs, rootDirRel, opts.AllowedGlobs)
		if err := autofixFunc(result, autofixOpts)
			err != nil {
			_, _ = fmt.Fprintln(procInout.Stderr, err.Error())
			return cli.ExitAbnormal
		}

		return cli.ExitNormal
	}
}
