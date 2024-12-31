package rpc

import (
	"testing"

	"github.com/anthony-dong/golang/pkg/idl"
)

func TestNewThriftServer(t *testing.T) {
	addr := ":10010"
	desc, err := idl.ParseThriftIDLFromProvider(idl.NewMemoryIDLProvider(GetTestIDLPath()))
	if err != nil {
		t.Fatal(err)
	}
	server := NewThriftServer(desc, addr)
	defer server.Close()
	go server.Run()
	runTestClient(t, "a.b.c", addr)
}
