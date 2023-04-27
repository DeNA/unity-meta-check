package chanutil

func ToSlice[T interface{}](ch <-chan T) []T {
	result := make([]T, 0)
	for x := range ch {
		result = append(result, x)
	}
	return result
}
