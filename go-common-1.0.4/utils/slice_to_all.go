package utils

func SliceToMap[K comparable, T any](list []T, f func(T) K) map[K]T {
	m := make(map[K]T)
	for _, t := range list {
		m[f(t)] = t
	}
	return m
}
func SliceToMapO[K comparable, T, D any](list []T, f func(T) (K, D)) map[K]D {
	m := make(map[K]D)
	for _, t := range list {
		k, v := f(t)
		m[k] = v
	}
	return m
}

func SliceToMapArray[K comparable, T any](list []T, f func(T) K) map[K][]T {
	m := make(map[K][]T)
	for _, t := range list {
		m[f(t)] = append(m[f(t)], t)
	}
	return m
}
func SliceToMapArrayO[K comparable, T, D any](list []T, f func(T) (K, D)) map[K][]D {
	m := make(map[K][]D)
	for _, t := range list {
		k, v := f(t)
		m[k] = append(m[k], v)
	}
	return m
}

func SliceToSlice[S, D any](list []S, f func(S) D, filter ...func(S) bool) []D {
	ids := make([]D, 0, len(list))

	if len(filter) > 0 {
		filterFunc := filter[0]
		for _, t := range list {
			if filterFunc(t) {
				ids = append(ids, f(t))
			}
		}
	} else {
		for _, t := range list {

			ids = append(ids, f(t))
		}
	}

	return ids
}
func SliceMerge[S any](list [][]S) []S {
	size := 0
	for _, t := range list {
		size += len(t)

	}
	rs := make([]S, 0, size)
	for _, d := range list {
		rs = append(rs, d...)
	}
	return rs
}
func CopyMaps[K comparable, T any](maps map[K]T) map[K]T {

	temp := make(map[K]T, len(maps))
	for k, t := range maps {
		temp[k] = t
	}

	return temp
}
