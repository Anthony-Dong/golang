package cpp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
)

// Tools
/**
Compile: clang++ -Wall -std=c++17 -O0 -g -I/usr/local/include -c ./cpp/fmt.cpp -o output/fmt.a
Link: clang++ -o output/fmt output/fmt.a -L/usr/local/lib -lspdlog -lgtest_main -lgtest
Run: output/fmt
*/

const ToolsOutputDir = "output"

type Tools struct {
	CXX string `json:",omitempty"`
	CC  string `json:",omitempty"`
	Dir string `json:",omitempty"`

	BuildArgs []string `json:",omitempty"` // 编译命令
	LinkArgs  []string `json:",omitempty"` // 链接命令

	SRCS []string `json:",omitempty"` // 源文件
	HDRS []string `json:",omitempty"` // 头文件

	BuildIncludes []string `json:",omitempty"`
	LinkIncludes  []string `json:",omitempty"`
	LinkLibraries []string `json:",omitempty"`

	CompileType string `json:",omitempty"` // 构建类型

	objects []string
}

const CompileTypeDebug = "debug"
const CompileTypeRelease = "release"
const LinkTypeBinary = "binary"
const LinkTypeLibrary = "library"

func (t *Tools) Compile(thread int) error {
	wg := errgroup.Group{}
	wg.SetLimit(thread)
	lock := sync.Mutex{}
	for _, _src := range t.SRCS {
		src := _src
		wg.Go(func() error {
			if obj, err := t.compile(src); err != nil {
				return err
			} else {
				lock.Lock()
				t.objects = append(t.objects, obj)
				lock.Unlock()
			}
			return nil
		})
	}
	return wg.Wait()
}

func (t *Tools) compile(src string) (string, error) {
	args := t.BuildArgs
	if !t.hasArgs(args, "-W") {
		args = append(args, "-Wall")
	}
	if !t.hasArgs(args, "-std") {
		args = append(args, fmt.Sprintf("-std=c++%s", CXX_STANDARD()))
	}
	switch t.CompileType {
	case CompileTypeRelease:
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
	object := t.NewObjectName(src)
	args = append(args, "-c", src)
	args = append(args, "-o", t.NewObjectName(src))
	command := exec.Command(t.CXX, args...)
	command.Dir = t.Dir
	logs.Info("Compile: %s", utils.PrettyCmd(command))
	if err := utils.RunCommand(command); err != nil {
		return "", err
	}
	return object, nil
}

func (t *Tools) NewObjectName(filename string) string {
	filename = filepath.Base(filename)
	return filepath.Join(ToolsOutputDir, strings.TrimSuffix(filename, filepath.Ext(filename))+".o")
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

func (t *Tools) Link(linkType string, file string) error {
	output := filepath.Join(ToolsOutputDir, file)
	switch linkType {
	case LinkTypeBinary:
		return t.linkBinary(output)
	case LinkTypeLibrary:
		if !strings.HasSuffix(output, ".a") {
			output = output + ".a"
		}
		return t.linkLibrary(output)
	}
	return fmt.Errorf(`not support link type`)
}

func (t *Tools) linkBinary(output string) error {
	args := []string{"-o", output}
	for _, object := range t.objects {
		args = append(args, object)
	}
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

func (t *Tools) linkLibrary(output string) error {
	// https://blog.csdn.net/xuhongning/article/details/6365200
	args := []string{"-r", "-c"}
	if logs.IsDebug() {
		args = append(args, "-v")
	}
	args = append(args, output)
	for _, object := range t.objects {
		args = append(args, object)
	}
	command := exec.Command("ar", args...)
	command.Dir = t.Dir
	logs.Info("Link: %s", utils.PrettyCmd(command))
	return utils.RunCommand(command)
}

func (t *Tools) Run(binaryName string) error {
	binaryFile := filepath.Join(ToolsOutputDir, binaryName)
	runCmd := exec.Command(binaryFile)
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

func CXX_STANDARD() string {
	if cc := os.Getenv("CXX_STANDARD"); cc != "" {
		return ""
	}
	return "17"
}
