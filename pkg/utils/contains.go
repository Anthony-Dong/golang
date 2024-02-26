package utils

func Contains[T comparable](arr []T, data T) bool {
	for _, elem := range arr {
		if elem == data {
			return true
		}
	}
	return false
}
