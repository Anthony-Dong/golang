package utils

import "github.com/mohae/deepcopy"

func DeepCopy[T any](input T) T {
	return deepcopy.Copy(input).(T)
}
