package pathchan

import "github.com/DeNA/unity-meta-check/util/typedpath"

func FromSlice(ss []typedpath.SlashPath) <-chan typedpath.SlashPath {
	ch := make(chan typedpath.SlashPath)
	go func() {
		defer close(ch)
		for _, str := range ss {
			ch <- str
		}
	}()
	return ch
}
