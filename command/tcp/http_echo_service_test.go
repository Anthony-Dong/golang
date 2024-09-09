package tcp

import (
	"context"
	"testing"
)

func Test_newHTTPEchoService(t *testing.T) {
	if err := newHTTPEchoService(context.Background(), ":8080"); err != nil {
		t.Fatal(err)
	}
}
