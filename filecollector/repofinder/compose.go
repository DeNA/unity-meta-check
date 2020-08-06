package repofinder

import (
	"github.com/DeNA/unity-meta-check/util/errutil"
	"sync"
)

func Compose(repoFinders []RepoFinder) RepoFinder {
	return func(writer chan<- *FoundRepo) error {
		var wg sync.WaitGroup
		errs := make([]error, 0)
		var errsMu sync.Mutex

		for _, findRepo := range repoFinders {
			wg.Add(1)
			go func(findRepo RepoFinder) {
				defer wg.Done()

				if err := findRepo(writer); err != nil {
					errsMu.Lock()
					errs = append(errs, err)
					errsMu.Unlock()
					return
				}
			}(findRepo)
		}

		wg.Wait()
		if len(errs) > 0 {
			return errutil.NewErrors(errs)
		}
		return nil
	}
}
