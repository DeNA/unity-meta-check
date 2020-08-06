package autofix

import (
	"fmt"
	"github.com/DeNA/unity-meta-check/unity"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"os"
)

type MetaRemover func(danglingMeta typedpath.RawPath) error

func NewMetaRemover(dryRun bool) MetaRemover {
	return func(danglingMeta typedpath.RawPath) error {
		stat, err := os.Stat(string(danglingMeta))
		if err != nil {
			return err
		}

		if !unity.IsMeta(danglingMeta.ToSlash()) || stat.IsDir() {
			return fmt.Errorf("must be a meta file: %s", danglingMeta)
		}

		if dryRun {
			return nil
		}
		return os.Remove(string(danglingMeta))
	}
}
