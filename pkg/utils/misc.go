package utils

import (
	"strings"
)

func ReadKV(str string) (string, string) {
	kvs := strings.SplitN(str, "=", 2)
	if len(kvs) != 2 {
		return "", ""
	}
	return strings.TrimSpace(kvs[0]), strings.TrimSpace(kvs[1])
}

func ReadKVBySep(str string, sep string) (string, string) {
	kvs := strings.SplitN(str, sep, 2)
	if len(kvs) != 2 {
		return "", ""
	}
	return strings.TrimSpace(kvs[0]), strings.TrimSpace(kvs[1])
}

func SplitKV(str string, sep string) (string, string) {
	return ReadKVBySep(str, sep)
}
