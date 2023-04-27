package chanutil

func FromSlice[T interface{}](ss []T) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, str := range ss {
			ch <- str
		}
	}()
	return ch
}
