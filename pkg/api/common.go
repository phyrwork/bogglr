package api

func MapOf[T any, U any](a []T, f func(T) U) []U {
	b := make([]U, len(a))
	for i, e := range a {
		b[i] = f(e)
	}
	return b
}

func PointersOf[T any](a []T) []*T {
	b := make([]*T, len(a))
	for i := range a {
		b[i] = &a[i]
	}
	return b
}

func MapPointersOf[T any, U any](a []T, f func(T) U) []*U {
	b := MapOf(a, f)
	return PointersOf(b)
}
