package tcpdump

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpMessage_DumpRequest(t *testing.T) {
	request, err := http.NewRequest("GET", "/api/v1", bytes.NewBuffer([]byte(`hello world`)))
	if err != nil {
		t.Fatal(err)
	}
	dumpRequest, err := httputil.DumpRequest(request, true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(dumpRequest))

	dumpRequest2, err := DumpRequest(request)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(dumpRequest2))
	assert.Equal(t, string(dumpRequest), string(dumpRequest2))
}

// 模拟一个简单的 HTTP 处理函数
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

func TestHttpMessage_DumpResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	dumpResp, err := httputil.DumpResponse(resp, true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(dumpResp))

	dumpResp2, err := DumpResponse(resp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(dumpResp2))
	assert.Equal(t, string(dumpResp), string(dumpResp2))
}
