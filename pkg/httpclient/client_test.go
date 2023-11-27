package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anthony-dong/golang/pkg/logs"
)

func TestHostClient_Get(t *testing.T) {
	defer logs.Flush()
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, request.Method, http.MethodGet)
		assert.Equal(t, request.URL.Path, "/api/v1/get")
		assert.Equal(t, request.Header.Get("Cookie"), "xxx")
		assert.Equal(t, request.URL.Query().Get("q"), "value")
		if _, err := writer.Write([]byte("hello world")); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()
	resp := ""
	client := HostClient{Host: server.URL, Auth: NewCookieAuth("xxx")}
	if err := client.Get(context.Background(), "/api/v1/get", map[string]string{"q": "value"}, &resp); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resp, "hello world")
}
