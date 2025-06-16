package main

import (
	"context"

	"github.com/anthony-dong/golang/command/devtool"

	"github.com/spf13/cobra"
)

func main() {
	devtool.Run(nil, func(ctx context.Context) (*cobra.Command, error) {
		return devtool.NewCommand(nil)
	})
}
