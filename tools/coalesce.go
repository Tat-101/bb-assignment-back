package tools

func Coalesce[T comparable](val, fallback T) T {
	var zero T
	if val == zero {
		return fallback
	}
	return val
}
