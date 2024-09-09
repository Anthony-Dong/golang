package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIP(t *testing.T) {
	ip, err := GetIP(true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ip)
}

func TestGetPort(t *testing.T) {
	{
		port, _ := GetPort(":8080")
		assert.Equal(t, port, 8080)
	}
	{
		port, _ := GetPort("[fe80::1]:8080")
		assert.Equal(t, port, 8080)
	}
}

func TestGetAllIP(t *testing.T) {
	fmt.Println(GetAllIP(true))
}
