package checker

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

func StubChecker(result *CheckResult, err error) Checker {
	return func(_ typedpath.RawPath, _ *Options) (*CheckResult, error) {
		return result, err
	}
}
