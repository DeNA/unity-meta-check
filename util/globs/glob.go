package globs

import (
	"github.com/DeNA/unity-meta-check/util/pathutil"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"path"
)

type Glob string

func MatchAny(p typedpath.SlashPath, globs []Glob) (bool, Glob, error) {
	ancestors := pathutil.AllAncestorsAndSelf(p)
	for _, ancestor := range ancestors {
		for _, glob := range globs {
			ok, err := path.Match(string(glob), string(ancestor))
			if err != nil {
				return false, "", err
			}
			if ok {
				return true, glob, nil
			}
		}
	}
	return false, "", nil
}
