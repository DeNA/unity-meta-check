package repofinder

func Const(results []*FoundRepo, err error) RepoFinder {
	return func(writer chan<- *FoundRepo) error {
		for _, found := range results {
			writer <- found
		}
		return err
	}
}

var Empty = Const(nil, nil)
