package pathchan

import "github.com/DeNA/unity-meta-check/util/typedpath"

func ToSlice(ch <-chan typedpath.SlashPath) []typedpath.SlashPath {
	result := make([]typedpath.SlashPath, 0)
	for x := range ch {
		result = append(result, x)
	}
	return result
}
