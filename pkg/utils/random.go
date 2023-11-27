package utils

import (
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
