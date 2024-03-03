package utils

import "sort"

func MapFromSlice[Input, Output any](input []Input, handler func(Input) Output) []Output {
	ret := make([]Output, 0, len(input))
	for _, elem := range input {
		ret = append(ret, handler(elem))
	}
	return ret
}

func MapFromMap[IK comparable, IV any, Output any](input map[IK]IV, handler func(key IK, value IV) Output) []Output {
	ret := make([]Output, 0, len(input))
	for k, v := range input {
		ret = append(ret, handler(k, v))
	}
	return ret
}

func FlatMapFromMap[IK comparable, IV any, Output any](input map[IK]IV, handler func(key IK, value IV) []Output) []Output {
	ret := make([]Output, 0, len(input))
	for k, v := range input {
		ret = append(ret, handler(k, v)...)
	}
	return ret
}

type KV[K comparable, V any] struct {
	Key   K
	Value V
}

func SortMap[K comparable, V any](input map[K]V, less func(i, j K) bool) []KV[K, V] {
	result := make([]KV[K, V], 0, len(input))
	for k, v := range input {
		result = append(result, KV[K, V]{Key: k, Value: v})
	}
	sort.Slice(result, func(i, j int) bool {
		return less(result[i].Key, result[j].Key)
	})
	return result
}
