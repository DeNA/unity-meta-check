package autofix

import (
	"errors"
	"github.com/DeNA/unity-meta-check/unity/checker"
)

func StubAutoFixer(skipped *checker.CheckResult, err error) AutoFixer {
	return func(_ *checker.CheckResult, _ *Options) (*checker.CheckResult, error) {
		return skipped, err
	}
}

type AutoFixerCallArgs struct {
	Result  *checker.CheckResult
	Options *Options
}

func SpyAutoFixer(inherited AutoFixer, callArgs *[]AutoFixerCallArgs) AutoFixer {
	if inherited == nil {
		inherited = StubAutoFixer(nil, errors.New("SPY_AUTO_FIXER"))
	}
	return func(result *checker.CheckResult, opts *Options) (*checker.CheckResult, error) {
		*callArgs = append(*callArgs, AutoFixerCallArgs{
			Result:  result,
			Options: opts,
		})
		return inherited(result, opts)
	}
}
