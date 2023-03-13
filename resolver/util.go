package resolver

func toPtr[T any](t T) *T {
	return &t
}
