package golang

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/utils"
)

func NewGoEnvCommand() (*cobra.Command, error) {
	version := ""
	cmd := &cobra.Command{
		Use: "env",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !isVersion(version) {
				return fmt.Errorf(`invalid version: %s`, version)
			}
			output := filepath.Join(utils.GetUserHomeDir(), "software", "gotemp")
			gohome, err := install(version, output)
			if err != nil {
				return err
			}
			fmt.Printf("export GOROOT=%s\n", gohome)
			fmt.Printf("export PATH=%s/bin:${PATH}\n", gohome)
			fmt.Printf("go version\n")
			return nil
		},
	}
	cmd.Flags().StringVar(&version, "version", "", "the version of go")
	return cmd, nil
}

// https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
func install(version string, output string) (string, error) {
	gohome := filepath.Join(output, "go"+version)
	if utils.ExistDir(gohome) {
		return gohome, nil
	}
	filename := fmt.Sprintf("go%s.%s-%s.tar.gz", version, runtime.GOOS, runtime.GOARCH)
	url := fmt.Sprintf(`https://go.dev/dl/%s`, filename)
	if !utils.ExistFile(filepath.Join(output, filename)) {
		if err := runCommand(exec.Command("wget", "-O", filename, url), output); err != nil {
			return "", err
		}
	}
	if err := runCommand(exec.Command("tar", "-zxvf", filename), output); err != nil {
		return "", err
	}
	if err := runCommand(exec.Command("mv", "go", "go"+version), output); err != nil {
		return "", err
	}
	return gohome, nil
}

func runCommand(command *exec.Cmd, pwd string) error {
	fmt.Println(strings.Join(command.Args, " "))
	command.Dir = pwd
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func isVersion(version string) bool {
	split := strings.Split(version, ".")
	for _, elem := range split {
		if _, err := strconv.ParseInt(elem, 10, 64); err != nil {
			return false
		}
	}
	return len(split) == 2 || len(split) == 3
}
