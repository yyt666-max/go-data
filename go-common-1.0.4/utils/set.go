package utils

type Set[T comparable] interface {
	Has(t T) bool
	Set(t ...T)
	ToList() []T
	Size() int
	Remove(t ...T)
}

type setMap[T comparable] map[T]struct{}

func NewSet[T comparable](vs ...T) Set[T] {
	s := make(setMap[T])

	for _, t := range vs {
		s[t] = struct{}{}
	}
	return s
}

func (s setMap[T]) Remove(t ...T) {
	for _, v := range t {
		delete(s, v)
	}
}
func (s setMap[T]) Has(t T) bool {
	_, h := s[t]
	return h
}

func (s setMap[T]) Set(ts ...T) {
	for _, t := range ts {
		s[t] = struct{}{}
	}

}

func (s setMap[T]) ToList() []T {
	l := make([]T, 0, len(s))
	for k := range s {
		l = append(l, k)
	}
	return l
}
func (s setMap[T]) Size() int {
	return len(s)
}
func SliceToSet[T comparable](list []T) Set[T] {
	s := make(setMap[T])

	for _, t := range list {
		s[t] = struct{}{}
	}
	return s
}

func Intersection[T comparable](a []T, b []T) []T {
	mb := make(map[T]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []T
	for _, x := range a {
		if _, found := mb[x]; found {
			diff = append(diff, x)
		}
	}
	return diff
}
