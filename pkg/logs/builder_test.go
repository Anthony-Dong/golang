package logs

import (
	"testing"
)

func TestBuilder(t *testing.T) {
	Builder().Info().Prefix("http").KV("method", "GET").KV("path", "/api/v1").Emit(nil)
}
