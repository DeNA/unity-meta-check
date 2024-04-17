package cmd

import (
	"flag"
	"fmt"
	"github.com/DeNA/unity-meta-check/report"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-junit/junit"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-junit/options"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/version"
	"io"
)

func NewMain() cli.Command {
	return func(args []string, procInout cli.ProcessInout, env cli.Env) cli.ExitStatus {
		opts, err := options.BuildOptions(args, procInout)
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
		result := parse(io.TeeReader(procInout.Stdin, procInout.Stdout))

		if err := junit.WriteToFile(result, opts.OutPath); err != nil {
			_, _ = fmt.Fprintln(procInout.Stderr, err.Error())
			return cli.ExitAbnormal
		}

		if !result.Empty() {
			return cli.ExitAbnormal
		}
		return cli.ExitNormal
	}
}
