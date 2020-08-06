package pathutil

import (
	"github.com/DeNA/unity-meta-check/util/typedpath"
	"path"
	"strings"
)

func SplitPathElements(targetPath typedpath.SlashPath) []typedpath.BaseName {
	result := make([]typedpath.BaseName, 0)
	dir := targetPath
	for dir != "" {
		newDir, file := path.Split(strings.TrimRight(string(dir), "/"))
		result = append(result, typedpath.BaseName(file))
		dir = typedpath.SlashPath(newDir)
	}
	reverse(result)
	return result
}

func reverse(baseNames []typedpath.BaseName) {
	for i := 0; i < len(baseNames)/2; i++ {
		j := len(baseNames) - i - 1
		baseNames[i], baseNames[j] = baseNames[j], baseNames[i]
	}
}
