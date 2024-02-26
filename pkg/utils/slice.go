package utils

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
