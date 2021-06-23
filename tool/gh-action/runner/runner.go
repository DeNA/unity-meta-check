package runner

import (
	"github.com/DeNA/unity-meta-check/report"
	"github.com/DeNA/unity-meta-check/resultfilter"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-autofix/autofix"
	prcomment "github.com/DeNA/unity-meta-check/tool/unity-meta-check-github-pr-comment/github"
	"github.com/DeNA/unity-meta-check/tool/unity-meta-check-junit/junit"
	"github.com/DeNA/unity-meta-check/unity/checker"
	"io"
	"time"
)

type Runner func(opts *Options) (bool, error)

func NewRunner(
	check checker.Checker,
	filter resultfilter.Filter,
	send prcomment.SendFunc,
	doAutofix autofix.AutoFixer,
	w io.Writer,
) Runner {
	return func(opts *Options) (bool, error) {
		startTime := time.Now()

		resultNotFiltered, err := check(opts.RootDirAbs, opts.CheckerOpts)
		if err != nil {
			return false, err
		}

		resultFiltered, err := filter(resultNotFiltered, opts.FilterOpts)
		if err != nil {
			return false, err
		}

		if err := report.WriteResult(w, resultFiltered); err != nil {
			return false, err
		}

		if opts.EnableJUnit {
			if err := junit.WriteToFile(resultFiltered, startTime, opts.JUnitOutPath); err != nil {
				return false, err
			}
		}

		if opts.EnablePRComment {
			if err := send(resultFiltered, opts.PRCommentOpts); err != nil {
				return false, err
			}
		}

		if opts.EnableAutofix {
			if err := doAutofix(resultFiltered, opts.AutofixOpts); err != nil {
				return false, err
			}
		}

		return resultFiltered.Empty(), nil
	}
}
