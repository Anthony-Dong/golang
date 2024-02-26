package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MustParseAddr(t *testing.T) {
	assert.Equal(t, MustParseAddr("127.0.0.1:8888"), &SimpleAddr{Network: "tcp", Addr: "127.0.0.1:8888"})
	assert.Equal(t, MustParseAddr("tcp://127.0.0.1:8888"), &SimpleAddr{Network: "tcp", Addr: "127.0.0.1:8888"})
}
