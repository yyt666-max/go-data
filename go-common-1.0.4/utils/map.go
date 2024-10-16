package utils

func MapToSlice[K comparable, T any, D any](m map[K]T, f func(k K, t T) D) []D {
	r := make([]D, 0, len(m))
	for k, v := range m {
		r = append(r, f(k, v))
	}
	return r
}
func MapKeys[K comparable, T any](m map[K]T) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}
func MapToSliceNoKey[K comparable, T any](m map[K]T) []T {
	r := make([]T, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}
func MapChange[K comparable, T any, D any](m map[K]T, f func(T) D) map[K]D {
	nm := make(map[K]D)
	for k, v := range m {
		nm[k] = f(v)
	}
	return nm
}
