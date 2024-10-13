package utils


func SliceToMap[T comparable](slice []T) map[T]struct{} {
	m := make(map[T]struct{})
	for _, v := range slice {
		m[v] = struct{}{}
	}
	return m
}



