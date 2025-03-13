package utils

import (
	"fmt"
	"testing"
)

func TestPrettySize(t *testing.T) {
	fmt.Println(PrettySize(786390663000))
}
