package proxy

import (
	"testing"
)

func TestNewThriftHandler(t *testing.T) {
	p := NewProxy(":10086", NewThriftHandler("10.37.10.60:8888", ConsoleRecorder))
	if err := p.Run(); err != nil {
		t.Fatal(err)
	}
}
