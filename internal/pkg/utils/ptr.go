package utils

func Ptr[T any](val T) *T {
	return &val
}

func DePtr[T any](val *T) T {
	if val == nil {
		var zero T
		return zero
	}
	return *val
}
