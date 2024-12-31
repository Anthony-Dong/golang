package cpp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadToolConfigFromFile(t *testing.T) {
	tools := &Tools{}
	err := ReadToolConfigFromKV(map[string]interface{}{
		"CXX":       "/opt/homebrew/opt/llvm@14/bin/clang++",
		"CC":        "/opt/homebrew/opt/llvm@14/bin/clang++",
		"CXX@linux": "/usr/lib/llvm-14/bin/clang++",
		"CC@linux":  "/usr/lib/llvm-14/bin/clang",
	}, "linux", tools)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, tools.CXX, "/usr/lib/llvm-14/bin/clang++")
	assert.Equal(t, tools.CC, "/usr/lib/llvm-14/bin/clang")
}
