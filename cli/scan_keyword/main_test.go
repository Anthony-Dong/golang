package main

import (
	"encoding/base64"
	"testing"
)

func TestBase64(t *testing.T) {
	t.Log(base64.StdEncoding.EncodeToString([]byte("xxxx")))
}
