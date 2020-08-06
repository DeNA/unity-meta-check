package pathutil

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"path"
	"strings"
)

func AllAncestorsAndSelf(targetPath typedpath.SlashPath) []typedpath.SlashPath {
	result := make([]typedpath.SlashPath, 0)

	current := strings.TrimRight(string(targetPath), "/")
	for current != "" && current != "." {
		result = append(result, typedpath.SlashPath(current))
		current = path.Dir(current)
	}

	return result
}
