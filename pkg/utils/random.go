package utils

import (
	"math"
	"math/rand"
	"strings"
)

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(size int) string {
	builder := strings.Builder{}
	builder.Grow(size)
	for x := 0; x < size; x++ {
		builder.WriteRune(letterRunes[rand.Intn(len(letterRunes))])
	}
	return builder.String()
}

func RandBytes(size int) []byte {
	output := make([]byte, size)
	for x := 0; x < size; x++ {
		output[x] = byte(rand.Intn(math.MaxUint8))
	}
	return output
}
