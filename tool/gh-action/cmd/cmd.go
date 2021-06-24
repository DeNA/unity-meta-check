package cmd

import (
	"flag"
	"fmt"
	"github.com/DeNA/unity-meta-check/git"
	umc "github.com/DeNA/unity-meta-check/options"
	"github.com/DeNA/unity-meta-check/resultfilter"
	"github.com/DeNA/unity-meta-check/tool/gh-action/options"
	"github.com/DeNA/unity-meta-check/tool/gh-action/runner"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-autofix/autofix"
	prcomment "github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/github"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/l10n"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-junit/junit"
	"github.com/DeNA/unity-meta-check/unity"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/unity/meta"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/ostestable"
	"github.com/DeNA/unity-meta-check/version"
)

func Main(args []string, procInout cli.ProcessInout, env cli.Env) cli.ExitStatus {
	parse := options.NewParser(umc.NewRootDirValidator(ostestable.NewIsDir()))

	opts, err := parse(args, procInout, env)
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

	validate := runner.NewValidateFunc(
		umc.NewUnityProjectDetector(logger),
		umc.NewIgnoredGlobsBuilder(logger),
		autofix.NewOptionsBuilder(ostestable.NewGetwd()),
		l10n.ReadTemplateFile,
	)
	runnerOpts, err := validate(opts.RootDirAbs, opts.UnsafeInputs, opts.Token)
	if err != nil {
		logger.Error(err.Error())
		return cli.ExitAbnormal
	}

	// NOTE: dry run is not necessary on GitHub Actions.
	dryRun := false
	//goland:noinspection GoBoolExpressions
	run := runner.NewRunner(
		checker.NewChecker(
			checker.NewStrategySelector(unity.NewFindPackages(logger), git.NewLsFiles(logger), logger),
			logger,
		),
		resultfilter.NewFilter(logger),
		junit.WriteToFile,
		prcomment.NewSendFunc(prcomment.NewPullRequestCommentSender(prcomment.NewHttp(), logger)),
		autofix.NewAutoFixer(
			dryRun,
			autofix.NewMetaTypeDetector(ostestable.NewIsDir()),
			autofix.NewMetaCreator(dryRun, meta.RandomGUIDGenerator(), logger),
			autofix.NewMetaRemover(dryRun),
			logger,
		),
		procInout.Stdout,
	)
	ok, err := run(runnerOpts)
	if err != nil {
		logger.Error(err.Error())
		return cli.ExitAbnormal
	}

	if ok {
		return cli.ExitNormal
	}
	return cli.ExitAbnormal
}
