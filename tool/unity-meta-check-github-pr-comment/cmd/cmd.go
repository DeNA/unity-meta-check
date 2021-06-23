package cmd

import (
	"flag"
	"fmt"
	"github.com/DeNA/unity-meta-check/report"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/github"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/options"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/version"
	"io"
)

func NewMain() cli.Command {
	return func(args []string, procInout cli.ProcessInout, env cli.Env) cli.ExitStatus {
		opts, err := options.BuildOptions(args, procInout, env)
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

		logger := logging.NewLogger(opts.LogLevel, procInout.Stderr)

		send := github.NewSendFunc(
			github.NewPullRequestCommentSender(github.NewHttp(), logger),
		)

		parse := report.NewParser()
		result := parse(io.TeeReader(procInout.Stdin, procInout.Stdout))

		if err := send(result, &github.Options{
			Tmpl:          opts.Tmpl,
			SendIfSuccess: opts.SendIfSuccess,
			Owner:         opts.Owner,
			Repo:          opts.Repo,
			PullNumber:    opts.PullNumber,
		}); err != nil {
			logger.Error(err.Error())
			return cli.ExitAbnormal
		}

		if !result.Empty() {
			return cli.ExitAbnormal
		}
		return cli.ExitNormal
	}
}
