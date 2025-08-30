package utils

func GetSafeValue[T any](value any) T {
	if value == nil {
		var emptyVal T

		return emptyVal
	}

	v, ok := value.(T)
	if !ok {
		var emptyVal T

		return emptyVal
	}

	return v
}
