package globs

import (
	"github.com/DeNA/unity-meta-check/util/pathutil"
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"path"
)

type Glob string

func (g Glob) ToSlash() typedpath.SlashPath {
	return typedpath.NewSlashPathUnsafe(string(g))
}

func MatchAny(p typedpath.SlashPath, globs []Glob, cwd typedpath.SlashPath) (bool, Glob, error) {
	pAbs := joinCwdIfRel(cwd, p)
	ancestors := pathutil.AllAncestorsAndSelf(pAbs)
	for _, ancestor := range ancestors {
		for _, relGlob := range globs {
			glob := joinCwdIfRel(cwd, relGlob.ToSlash())
			ok, err := path.Match(string(glob), string(ancestor))
			if err != nil {
				return false, "", err
			}
			if ok {
				return true, relGlob, nil
			}
		}
	}
	return false, "", nil
}

func joinCwdIfRel(cwd typedpath.SlashPath, path typedpath.SlashPath) typedpath.SlashPath {
	if path.IsAbs() {
		return path
	}
	return cwd.JoinSlashPath(path)
}