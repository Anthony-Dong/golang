package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Printf("args: %#v\n", os.Args)
	for _, elem := range os.Environ() {
		fmt.Printf("env: %q\n", elem)
	}
	time.Sleep(time.Second * 10000)
}
