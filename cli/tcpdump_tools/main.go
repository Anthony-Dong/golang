package main

import (
	"context"
	"log"
	"os"

	"github.com/anthony-dong/golang/cli/tcpdump_tools/tcpdump"
)

func main() {
	command, err := tcpdump.NewCommand()
	if err != nil {
		exitError(err)
	}
	if err := command.ExecuteContext(context.Background()); err != nil {
		exitError(err)
	}
}

func exitError(err error) {
	log.Printf("[%s] %s\n", "tcpdump", err.Error())
	os.Exit(1)
}
