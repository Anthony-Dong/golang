package golang

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/anthony-dong/golang/command"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/tools"
	"github.com/anthony-dong/golang/pkg/utils"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "golang",
		Short: "Go language related tools and utilities",
	}
	if err := command.AddCommand(cmd, NewGoRunCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, NewGoTestCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, NewListGoModSize); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, NewGoEnvCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, NewGolandInitCommand); err != nil {
		return nil, err
	}
	return cmd, nil
}

func NewGolandInitCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "goland-init",
		Short: "Initialize Goland project dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := RunCommand("go", "mod", "tidy"); err != nil {
				return err
			}

			// go list -m -json -mod=mod all
			return RunCommand("go", "list", "-m", "-json", "-mod=mod", "all")
		},
	}
	return cmd, nil
}

func RunCommand(args ...string) error {
	log.Printf("exec: %s\n", strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "run command %v error", args)
	}
	return nil
}

func NewGoRunCommand() (*cobra.Command, error) {
	runPkg := ""
	isDebug := false
	runEnv := make([]string, 0)
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run Go programs",
		RunE: func(cmd *cobra.Command, args []string) error {
			var command *exec.Cmd
			if isDebug {
				dlv := tools.NewDebugDlv(runPkg, args...)
				command = exec.Command(dlv.Name(), dlv.Args()...)
			} else {
				command = exec.Command("go", append([]string{"run", "-v", runPkg}, args...)...)
			}
			command.Env = append(os.Environ(), runEnv...)
			return utils.RunCmd(command, "")
		},
	}
	cmd.Flags().BoolVar(&isDebug, "debug", false, "enable debug")
	cmd.Flags().StringVar(&runPkg, "run", ".", "go run pkg name")
	cmd.Flags().StringSliceVar(&runEnv, "env", []string{}, "go test env")
	if err := cmd.MarkFlagRequired("run"); err != nil {
		return nil, err
	}
	return cmd, nil
}

func NewGoTestCommand() (*cobra.Command, error) {
	testName := ""
	testPkg := ""
	testEnv := make([]string, 0)
	isDebug := false
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Run Go tests",
		RunE: func(cmd *cobra.Command, args []string) error {
			var command *exec.Cmd
			if isDebug {
				dlv := tools.NewTestDlv(testName, testPkg)
				command = exec.Command(dlv.Name(), dlv.Args()...)
			} else {
				command = exec.Command("go", "test", "-v", fmt.Sprintf("-run=%s", testName), "-count=1", testPkg)
			}
			command.Env = append(os.Environ(), testEnv...)
			return utils.RunCmd(command, "")
		},
	}

	cmd.Flags().BoolVar(&isDebug, "debug", false, "enable debug")
	cmd.Flags().StringVar(&testName, "run", "", "go test name")
	cmd.Flags().StringVar(&testPkg, "pkg", "", "go test pkg")
	cmd.Flags().StringSliceVar(&testEnv, "env", []string{}, "go test env")
	if err := cmd.MarkFlagRequired("run"); err != nil {
		return nil, err
	}
	if err := cmd.MarkFlagRequired("pkg"); err != nil {
		return nil, err
	}
	return cmd, nil
}
