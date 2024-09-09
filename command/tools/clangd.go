package tools

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Command struct {
	File      string `json:"file"`
	Command   string `json:"command"`
	Directory string `json:"directory"`
}

func NewClangdCommand() (*cobra.Command, error) {
	filename := ""
	action := ""
	command := cobra.Command{
		Use:   "clangd",
		Short: "parse compile_commands.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			if filename == "" {
				return fmt.Errorf(`the compile_commands.json connot be empty`)
			}
			file, err := os.Open(filename)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			decoder := json.NewDecoder(file)
			// 读取开始的 [
			if _, err := decoder.Token(); err != nil {
				return err
			}
			// while the array contains values
			for decoder.More() {
				var command Command
				err := decoder.Decode(&command)
				if err != nil {
					return err
				}
				if action == "list" {
					fmt.Println(command.File)
				}
			}
			// 结束读取 ]
			if _, err := decoder.Token(); err != nil {
				return err
			}
			return nil
		},
	}
	command.Flags().StringVarP(&filename, "read", "r", "compile_commands.json", "read the compile_commands.json")
	command.Flags().StringVarP(&action, "action", "q", "list", "action. eg: list")
	return &command, nil
}
