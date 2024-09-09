package main

import (
	"fmt"
	"os"
)

func main() {
	for _, elem := range os.Environ() {
		fmt.Println(elem)
	}
}
