package utils

import (
	"encoding/base64"
	"testing"
)

func TestRandBytes(t *testing.T) {
	bytes := RandBytes(1024 * 10)
	t.Log(base64.StdEncoding.EncodeToString(bytes))
}
