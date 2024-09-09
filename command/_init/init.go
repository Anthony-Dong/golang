package _init

import (
	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{Use: "init", Short: `quickly initialize various types of projects`}
	if err := command.AddCommand(cmd, NewCmakeCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, NewVscodeLaunch); err != nil {
		return nil, err
	}
	return cmd, nil
}
