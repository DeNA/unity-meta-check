package main

import (
	"flag"
	"fmt"
	"github.com/DeNA/unity-meta-check/git"
	"github.com/DeNA/unity-meta-check/options"
	"github.com/DeNA/unity-meta-check/report"
	"github.com/DeNA/unity-meta-check/resultfilter"
	"github.com/DeNA/unity-meta-check/unity"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/cli"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/version"
	"os"
)

func main() {
	var cmd cli.Command
	profModeSwitch := os.Getenv("UNITY_META_CHECK_PROFILE")
	if profModeSwitch == "cpu" {
		cmd = cli.NewCommandWithCPUProfile(NewMain())
	} else if profModeSwitch == "heap" {
		cmd = cli.NewCommandWithHeapProfile(NewMain())
	} else if profModeSwitch != "" {
		println(fmt.Sprintf("unsupported profile mode: %q", profModeSwitch))
		os.Exit(1)
	} else {
		cmd = NewMain()
	}

	exitStatus := cmd(os.Args[1:], cli.GetProcessInout(), cli.NewEnv())
	os.Exit(int(exitStatus))
}

func NewMain() cli.Command {
	return func(args []string, procInout cli.ProcessInout, env cli.Env) cli.ExitStatus {
		opts, err := options.ParseArgs(args, procInout)
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

		check := checker.NewChecker(
			checker.NewStrategySelector(
				unity.NewFindPackages(logger),
				git.NewLsFiles(logger),
				logger,
			),
			logger,
		)
		result, err := check(
			opts.RootDirAbs,
			&checker.Options{
				IgnoreCase:                opts.IgnoreCase,
				IgnoreSubmodulesAndNested: opts.IgnoreSubmodulesAndNested,
				TargetType:                opts.TargetType,
			},
		)
		if err != nil {
			_, _ = fmt.Fprintln(procInout.Stderr, err.Error())
			return cli.ExitAbnormal
		}

		filterFunc := resultfilter.NewFilter(logger)
		filtered, err := filterFunc(result, &resultfilter.Options{
			IgnoreDangling: opts.IgnoreDangling,
			IgnoredGlobs:   opts.IgnoredPaths,
		})
		if err != nil {
			_, _ = fmt.Fprintln(procInout.Stderr, err.Error())
			return cli.ExitAbnormal
		}

		if err := report.WriteResult(procInout.Stdout, filtered); err != nil {
			_, _ = fmt.Fprintln(procInout.Stderr, err.Error())
			return cli.ExitAbnormal
		}

		if !filtered.Empty() {
			return cli.ExitAbnormal
		}
		return cli.ExitNormal
	}
}
