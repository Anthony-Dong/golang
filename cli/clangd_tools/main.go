package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Command struct {
	File      string `json:"file"`
	Command   string `json:"command"`
	Directory string `json:"directory"`
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	// 读取开始的 [
	if _, err := decoder.Token(); err != nil {
		log.Fatal(err)
	}

	// while the array contains values
	for decoder.More() {
		var command Command

		// decode an array value (Command)
		err := decoder.Decode(&command)
		if err != nil {
			panic(err)
		}

		// 打印或者处理command
		fmt.Println(command.File)
	}

	// 读取结束的 ]
	if _, err := decoder.Token(); err != nil {
		panic(err)
	}

}
