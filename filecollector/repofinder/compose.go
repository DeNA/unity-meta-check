package repofinder

func Compose(repoFinders []RepoFinder) RepoFinder {
	return func() ([]FoundRepo, error) {
		result := make([]FoundRepo, 0)
		for _, repoFinder := range repoFinders {
			foundRepo, err := repoFinder()
			if err != nil {
				return nil, err
			}
			result = append(result, foundRepo...)
		}
		return result, nil
	}
}
