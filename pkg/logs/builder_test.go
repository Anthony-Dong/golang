package logs

import (
	"testing"
)

func TestBuilder(t *testing.T) {
	Builder().Info().String("http:").KV("method", "GET").KV("path", "/api/v1").Emit(nil)
}
