package runner

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/report"
	"github.com/DeNA/unity-meta-check/resultfilter"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-autofix/autofix"
	prcomment "github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/github"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-junit/junit"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"github.com/DeNA/unity-meta-check/util/logging"
	"io"
	"time"
)

type Runner func(opts *Options) (bool, error)

func NewRunner(
	check checker.Checker,
	filter resultfilter.Filter,
	writeJunitXML junit.WriteToFileFunc,
	send prcomment.SendFunc,
	doAutofix autofix.AutoFixer,
	w io.Writer,
	logger logging.Logger,
) Runner {
	return func(opts *Options) (bool, error) {
		startTime := time.Now()

		logger.Debug(fmt.Sprintf("check: %#v", opts.CheckerOpts))
		resultNotFiltered, err := check(opts.RootDirAbs, opts.CheckerOpts)
		if err != nil {
			return false, err
		}

		logger.Debug(fmt.Sprintf("filter result: %#v", opts.FilterOpts))
		resultFiltered, err := filter(resultNotFiltered, opts.FilterOpts)
		if err != nil {
			return false, err
		}

		if opts.EnableAutofix {
			logger.Debug(fmt.Sprintf("autofix: %#v", opts.AutofixOpts))
			skipped, err := doAutofix(resultFiltered, opts.AutofixOpts)
			if err != nil {
				return false, err
			}
			resultFiltered = skipped
		} else {
			logger.Debug(`skip autofix because "enable_autofix" is false`)
		}

		logger.Debug("print check result")
		if err := report.WriteResult(w, resultFiltered); err != nil {
			return false, err
		}

		if opts.EnableJUnit {
			logger.Debug(fmt.Sprintf("write junit report: %q", opts.JUnitOutPath))
			if err := writeJunitXML(resultFiltered, startTime, opts.JUnitOutPath); err != nil {
				return false, err
			}
		} else {
			logger.Debug(`skip write junit report because "enable_junit" is false`)
		}

		if opts.EnablePRComment {
			logger.Debug(fmt.Sprintf("send pull request comment if necessary:\n%s", prcomment.MaskOptions(opts.PRCommentOpts)))
			if err := send(resultFiltered, opts.PRCommentOpts); err != nil {
				return false, err
			}
		} else {
			logger.Debug(`skip send a pull request comment because "enable_pr_comment" is false`)
		}

		return resultFiltered.Empty(), nil
	}
}
