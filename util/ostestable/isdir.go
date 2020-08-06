package ostestable

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"os"
)

type IsDir func(path typedpath.RawPath) (bool, error)

func NewIsDir() IsDir {
	return func(path typedpath.RawPath) (bool, error) {
		stat, err := os.Stat(string(path))
		if err != nil {
			return false, err
		}
		return stat.IsDir(), nil
	}
}
