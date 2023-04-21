package filecollector

import "github.com/DeNA/unity-meta-check/util/typedpath"

type Entry struct {
	Path  typedpath.SlashPath
	IsDir bool
}
