package filecollector

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

func StubFileAggregator(result []Entry, err error) FileAggregator {
	return func(rootDirAbs typedpath.RawPath, opts *Options, ch chan<- Entry) error {
		for _, path := range result {
			ch <- path
		}
		return err
	}
}

func StubSuccessfulFileAggregator(result []Entry) FileAggregator {
	return StubFileAggregator(result, nil)
}
