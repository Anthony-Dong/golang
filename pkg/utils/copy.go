package utils

func CopyMap(kv map[string]string) map[string]string {
	result := make(map[string]string, len(kv))
	for k, v := range kv {
		result[k] = v
	}
	return result
}
