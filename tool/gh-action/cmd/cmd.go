package cmd

import (
	"flag"
	"fmt"
	"github.com/DeNA/unity-meta-check/git"
	common "github.com/DeNA/unity-meta-check/options"
	"github.com/DeNA/unity-meta-check/resultfilter"
	"github.com/DeNA/unity-meta-check/tool/gh-action/inputs"
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
	"github.com/DeNA/unity-meta-check/util/testutil"
	"github.com/DeNA/unity-meta-check/version"
)

func Main(args []string, procInout cli.ProcessInout, env cli.Env) cli.ExitStatus {
	parse := options.NewParser()

	opts, err := parse(args, procInout, env)
	if err != nil {
		if err != flag.ErrHelp {
			_, _ = fmt.Fprintln(procInout.Stderr, err.Error())
		}
		return cli.ExitAbnormal
	}

	logger := logging.NewLogger(logging.MustParseSeverity(opts.Inputs.LogLevel), procInout.Stderr)
	logger.Debug(fmt.Sprintf("inputs=%#v", opts.Inputs))
	logger.Debug(inputs.MaskedActionEnv(opts.Env))

	if opts.Version {
		_, _ = fmt.Fprintln(procInout.Stdout, version.Version)
		return cli.ExitNormal
	}

	// XXX: Avoid errors like "detected dubious ownership in repository at '...'".
	// 	    The file owner of files in opts.Env.Workspace is user outside the container, but git is run by root in the container.
	// 	    Therefore, a wrong owner seem from git in a container. So git throw the error.
	// 	    This code trust all repositories in both the container and mounted volumes. It justified by the 2 assumptions:
	//
	//        1) No evil repositories in the container. It is unity-meta-check's responsibility.
	//        2) No evil repositories in mounted volumes. It is user's responsibility.
	//
	gitGlobalConfig := git.NewGlobalConfig(logger)
	if err := gitGlobalConfig(testutil.SpyWriteCloser(nil), "--add", "safe.directory", "*"); err != nil {
		logger.Error(fmt.Sprintf("failed to add '*' to safe.directory: %s", err.Error()))
		return cli.ExitAbnormal
	}

	validate := runner.NewValidateFunc(
		common.NewRootDirValidator(ostestable.NewIsDir()),
		common.NewUnityProjectDetector(logger),
		common.NewIgnoredGlobsBuilder(logger),
		autofix.NewOptionsBuilder(ostestable.NewGetwd()),
		l10n.ReadTemplateFile,
		inputs.NewReadEventPayload(logger),
	)
	runnerOpts, err := validate(opts.Inputs, opts.Env)
	if err != nil {
		logger.Error(err.Error())
		return cli.ExitAbnormal
	}
	logger.Debug(fmt.Sprintf("runner options: %#v", runnerOpts))

	// NOTE: dry run is not necessary on GitHub Actions.
	dryRun := false
	//goland:noinspection GoBoolExpressions
	run := runner.NewRunner(
		checker.NewChecker(
			checker.NewStrategySelector(unity.NewFindPackages(logger), git.NewLsFiles(logger), logger),
			logger,
		),
		resultfilter.NewFilter(ostestable.NewGetwd(), logger),
		junit.WriteToFile,
		prcomment.NewSendFunc(prcomment.NewPullRequestCommentSender(prcomment.NewHttp(), logger)),
		autofix.NewAutoFixer(
			dryRun,
			ostestable.NewGetwd(),
			autofix.NewMetaTypeDetector(ostestable.NewIsDir()),
			autofix.NewMetaCreator(dryRun, meta.RandomGUIDGenerator(), logger),
			autofix.NewMetaRemover(dryRun),
			logger,
		),
		procInout.Stdout,
		logger,
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
