package utils

import "strings"

func ReadKV(str string) (string, string) {
	kvs := strings.SplitN(str, "=", 2)
	if len(kvs) != 2 {
		return "", ""
	}
	return strings.TrimSpace(kvs[0]), strings.TrimSpace(kvs[1])
}
