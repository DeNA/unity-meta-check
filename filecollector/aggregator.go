package filecollector

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/filecollector/repofinder"
	"github.com/DeNA/unity-meta-check/git"
	"github.com/DeNA/unity-meta-check/util/errutil"
	"github.com/DeNA/unity-meta-check/util/logging"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"sync"
)

type FileAggregator func(rootDirAbs typedpath.RawPath, opts *Options, ch chan<- Entry) error

func NewFileAggregator(gitLsFiles git.LsFiles, findRepo repofinder.RepoFinder, logger logging.Logger) FileAggregator {
	collectFiles := New(gitLsFiles, logger)

	return func(rootDirAbs typedpath.RawPath, opts *Options, ch chan<- Entry) error {
		var errsMutex sync.Mutex
		errs := make([]error, 0)

		foundRepos, err := findRepo()
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		for _, foundRepo := range foundRepos {
			logger.Info(fmt.Sprintf("repository found: %q (submodule=%t)", foundRepo.RawPath, foundRepo.Type))

			wg.Add(1)
			go func(foundRepo repofinder.FoundRepo) {
				defer wg.Done()
				ch <- Entry{Path: foundRepo.RawPath.ToSlash(), IsDir: true}
			}(foundRepo)

			wg.Add(1)
			go func(foundRepo repofinder.FoundRepo) {
				defer wg.Done()
				if err := collectFiles(rootDirAbs, foundRepo.RawPath, opts, ch); err != nil {
					errsMutex.Lock()
					errs = append(errs, err)
					errsMutex.Unlock()
					return
				}
			}(foundRepo)
		}

		if err = collectFiles(rootDirAbs, ".", opts, ch); err != nil {
			errsMutex.Lock()
			errs = append(errs, err)
			errsMutex.Unlock()
		}
		wg.Wait()

		if len(errs) > 0 {
			return errutil.NewErrors(errs)
		}
		return nil
	}
}
