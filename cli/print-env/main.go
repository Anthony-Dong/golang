package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("args: %#v\n", os.Args)
	for _, elem := range os.Environ() {
		fmt.Printf("env: %q\n", elem)
	}
}
