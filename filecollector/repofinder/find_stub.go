package repofinder

func StubRepoFinder(result []FoundRepo, err error) RepoFinder {
	return func() ([]FoundRepo, error) {
		return result, err
	}
}
