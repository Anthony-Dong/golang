package proxy

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stripPort(t *testing.T) {
	assert.Equal(t, stripPort("www.baidu.com:443"), "www.baidu.com")
	assert.Equal(t, stripPort("localhost:80"), "localhost")
	assert.Equal(t, stripPort("127.0.0.1:80"), "127.0.0.1")
	assert.Equal(t, stripPort("[fe80::216:3eff:fe53:5fa]:80"), "fe80::216:3eff:fe53:5fa")
	assert.Equal(t, net.ParseIP("fe80::216:3eff:fe53:5fa").String(), "fe80::216:3eff:fe53:5fa")
}
