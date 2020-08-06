package meta

import (
	"io"
)

type Gen interface {
	WriteTo(writer io.Writer) (int64, error)
}

