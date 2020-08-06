package filecollector

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

func StubFileAggregator(result []typedpath.SlashPath, err error) FileAggregator {
	return func(rootDirAbs typedpath.RawPath, opts *Options, ch chan<- typedpath.SlashPath) error {
		for _, path := range result {
			ch <- path
		}
		return err
	}
}

func StubSuccessfulFileAggregator(result []typedpath.SlashPath) FileAggregator {
	return StubFileAggregator(result, nil)
}
