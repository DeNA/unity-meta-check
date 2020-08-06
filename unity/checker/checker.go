package checker

import (
	"github.com/DeNA/unity-meta-check/filecollector"
	"github.com/DeNA/unity-meta-check/git"
	"github.com/DeNA/unity-meta-check/unity"
	"github.com/DeNA/unity-meta-check/util/errutil"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"sync"
)

type Checker func(rootDirAbs typedpath.RawPath, opts *Options) (*CheckResult, error)

func Check(rootDirAbs typedpath.RawPath, opts *Options, logger logging.Logger) (*CheckResult, error) {
	selectStrategy := NewStrategySelector(unity.NewFindPackages(logger), git.NewLsFiles(logger), logger)
	strategy, err := selectStrategy(rootDirAbs, opts)
	if err != nil {
		return nil, err
	}

	check := newChecker(strategy, logger)
	return check(rootDirAbs, opts)
}

func newChecker(strategy Strategy, logger logging.Logger) Checker {
	return func(rootDirAbs typedpath.RawPath, opts *Options) (*CheckResult, error) {
		ch := make(chan typedpath.SlashPath)

		var wg sync.WaitGroup
		var errsMu sync.Mutex
		errs := make([]error, 0)

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(ch)
			if err := strategy.CollectFiles(rootDirAbs, &filecollector.Options{IgnoreCase: opts.IgnoreCase}, ch); err != nil {
				errsMu.Lock()
				errs = append(errs, err)
				errsMu.Unlock()
				return
			}
		}()

		check := NewCheckingWorker(strategy.RequiresMeta, logger)
		result, err := check(opts.IgnoreCase, ch)
		if err != nil {
			errsMu.Lock()
			errs = append(errs, err)
			errsMu.Unlock()
		}
		if len(errs) > 0 {
			return nil, errutil.NewErrors(errs)
		}
		return result, nil
	}
}
