package cmd

import (
	"flag"
	"fmt"
	options2 "github.com/DeNA/unity-meta-check/options"
	"github.com/DeNA/unity-meta-check/report"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-autofix/autofix"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-autofix/options"
	"github.com/DeNA/unity-meta-check/unity/meta"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/ostestable"
	"github.com/DeNA/unity-meta-check/version"
)

func NewMain() cli.Command {
	return func(args []string, procInout cli.ProcessInout, env cli.Env) cli.ExitStatus {
		parseOpts := options.NewParser(options2.NewRootDirValidator(ostestable.NewIsDir()))
		opts, err := parseOpts(args, procInout)
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

		autofixFunc := autofix.NewAutoFixer(
			opts.DryRun,
			autofix.NewMetaTypeDetector(ostestable.NewIsDir()),
			autofix.NewMetaCreator(opts.DryRun, meta.RandomGUIDGenerator(), logger),
			autofix.NewMetaRemover(opts.DryRun),
			logger,
		)

		buildOpts := autofix.NewOptionsBuilder(ostestable.NewGetwd())
		autofixOpts, err := buildOpts(opts.RootDirAbs, opts.AllowedGlobs)
		if err != nil {
			logger.Error(err.Error())
			return cli.ExitNormal
		}

		if err := autofixFunc(result, autofixOpts); err != nil {
			logger.Error(err.Error())
			return cli.ExitAbnormal
		}

		return cli.ExitNormal
	}
}
