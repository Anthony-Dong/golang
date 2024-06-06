package rpc

import "testing"

func TestRequestToString(t *testing.T) {
	req, err := NewRpcRequest("xxx.xxx.xxx/Test?k1=v1&k2=v2", []string{"h1:v1"}, `{
		"k1": "v1"
	}`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(req.String())
}
