package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatSize(t *testing.T) {
	assert.Equal(t, FormatSize(1), "1B")
	assert.Equal(t, FormatSize(1024), "1KB")
	assert.Equal(t, FormatSize(1024*1024), "1M")
	assert.Equal(t, FormatSize(1024*1024*1024), "1G")
	assert.Equal(t, FormatSize(1024*1024*1024*1024), "1024G")

	assert.Equal(t, FormatSize(481012), "469.74KB")
	assert.Equal(t, FormatSize(71581320), "68.27M")
}
