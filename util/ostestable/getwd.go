package ostestable

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
)

type Getwd func() (typedpath.RawPath, error)

func NewGetwd() Getwd {
	return typedpath.Getwd
}
