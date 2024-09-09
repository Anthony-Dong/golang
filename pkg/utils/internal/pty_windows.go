//go:build windows

package internal

import (
	"os"
	"os/exec"
)

func RunPty(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
