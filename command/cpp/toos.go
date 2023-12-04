package cpp

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
)

// Tools
/**
Build: clang++ -Wall -std=c++17 -O0 -g -I/usr/local/include -c ./cpp/fmt.cpp -o output/fmt.a
Link: clang++ -o output/fmt output/fmt.a -L/usr/local/lib -lspdlog -lgtest_main -lgtest
Run: output/fmt
*/
type Tools struct {
	CXX string
	CC  string
	Dir string

	BuildArgs []string
	LinkArgs  []string

	SRCS []string
	HDRS []string

	BuildIncludes []string
	LinkIncludes  []string
	LinkLibraries []string

	Binary string

	BuildType string
}

const BuildTypeRelease = "release"

func (t *Tools) Build() error {
	args := t.BuildArgs

	if !t.hasArgs(args, "-W") {
		args = append(args, "-Wall")
	}

	if !t.hasArgs(args, "-std") {
		args = append(args, "-std=c++17")
	}

	switch t.BuildType {
	case BuildTypeRelease:
		if t.hasArgs(args, "-O") {
			logs.Debug("Build: cannot set arg -O2 because you already configured it")
		} else {
			args = append(args, "-O2")
		}
	default:
		if t.hasArgs(args, "-O") {
			logs.Debug("Build: cannot set arg -O0 because you already configured it")
		} else {
			args = append(args, "-O0")
		}
		if t.hasArgs(args, "-O0") { // only build -O append -g args
			args = append(args, "-g") // Generate source-level debug information
		}
	}

	for _, elem := range t.BuildIncludes {
		args = append(args, fmt.Sprintf("-I%s", elem))
	}
	args = append(args, "-c")
	for _, header := range t.HDRS {
		args = append(args, header)
	}
	for _, src := range t.SRCS {
		args = append(args, src)
	}
	args = append(args, "-o", t.Binary+".a")
	command := exec.Command(t.CXX, args...)
	command.Dir = t.Dir
	logs.Info("Build: %s", utils.PrettyCmd(command))
	return utils.RunCommand(command)
}

func (t *Tools) Init() error {
	if t.CXX == "" {
		return fmt.Errorf(`not found cxx`)
	}
	if t.CC == "" {
		return fmt.Errorf(`not found cc`)
	}
	return nil
}

func (t *Tools) Link() error {
	args := []string{"-o", t.Binary}
	args = append(args, t.Binary+".a")
	for _, elem := range t.LinkIncludes {
		args = append(args, fmt.Sprintf("-L%s", elem))
	}
	for _, elem := range t.LinkLibraries {
		args = append(args, fmt.Sprintf("-l%s", elem))
	}
	command := exec.Command(t.CXX, args...)
	command.Dir = t.Dir
	logs.Info("Link: %s", utils.PrettyCmd(command))
	return utils.RunCommand(command)
}

func (t *Tools) Run() error {
	runCmd := exec.Command(t.Binary)
	runCmd.Dir = t.Dir
	logs.Info("Run: %s", utils.PrettyCmd(runCmd))
	return utils.RunCommand(runCmd)
}

func (t *Tools) hasArgs(args []string, arg string) bool {
	for _, elem := range args {
		if strings.HasPrefix(elem, arg) {
			return true
		}
	}
	return false
}

func CXX() string {
	if cxx := os.Getenv("CXX"); cxx != "" {
		return cxx
	}
	return "clang++"
}

func CC() string {
	if cc := os.Getenv("CC"); cc != "" {
		return cc
	}
	return "clang"
}
