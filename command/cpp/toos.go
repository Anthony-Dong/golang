package cpp

import (
	"context"
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

type Tools struct {
	CXX    string `json:"CXX,omitempty" yaml:"CXX,omitempty"`
	CC     string `json:"CC,omitempty" yaml:"CC,omitempty"`
	Pwd    string `json:"Pwd,omitempty" yaml:"Pwd,omitempty"`
	Output string `json:"Output,omitempty" yaml:"Output,omitempty"`

	BuildArgs []string `json:"BuildArgs,omitempty" yaml:"BuildArgs,omitempty"` // 编译命令
	LinkArgs  []string `json:"LinkArgs,omitempty" yaml:"LinkArgs,omitempty"`   // 链接命令

	SRCS      []string `json:"SRCS,omitempty" yaml:"SRCS,omitempty"`           // 源文件
	HDRS      []string `json:"HDRS,omitempty" yaml:"HDRS,omitempty"`           // 头文件
	BuildType string   `json:"BuildType,omitempty" yaml:"BuildType,omitempty"` // 构建类型

	objects []string
}

const BuildTypeDebug = "debug"
const BuildTypeRelease = "release"
const LinkTypeBinary = "binary"
const LinkTypeLibrary = "library"

func (t *Tools) Build(ctx context.Context, thread int) error {
	wg := errgroup.Group{}
	wg.SetLimit(thread)
	lock := sync.Mutex{}
	for _, _src := range t.SRCS {
		src := _src
		wg.Go(func() error {
			if obj, err := t.build(ctx, src); err != nil {
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

func (t *Tools) build(ctx context.Context, src string) (string, error) {
	args := t.BuildArgs
	if !t.hasArgs(args, "-W") {
		args = append(args, "-Wall")
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
	object := t.NewObjectName(src)
	args = append(args, "-c", src)
	args = append(args, "-o", object)
	command := exec.Command(t.CXX, args...)
	command.Dir = t.Pwd
	logs.CtxDebug(ctx, "Build: %s", utils.PrettyCmd(command))
	if err := utils.RunCommand(command); err != nil {
		return "", err
	}
	return object, nil
}

func (t *Tools) NewObjectName(filename string) string {
	filename = filepath.Base(filename)
	outputDir := filepath.Dir(t.Output)
	return filepath.Join(outputDir, strings.TrimSuffix(filename, filepath.Ext(filename))+".o")
}

func (t *Tools) Link(ctx context.Context, linkType string) error {
	switch linkType {
	case LinkTypeBinary:
		return t.linkBinary(ctx, t.Output)
	case LinkTypeLibrary:
		return t.linkLibrary(ctx, t.Output)
	}
	return fmt.Errorf(`not support link type`)
}

func (t *Tools) linkBinary(ctx context.Context, output string) error {
	args := []string{"-o", output}
	for _, object := range t.objects {
		args = append(args, object)
	}
	args = append(args, t.LinkArgs...)
	command := exec.Command(t.CXX, args...)
	command.Dir = t.Pwd
	logs.CtxDebug(ctx, "Link: %s", utils.PrettyCmd(command))
	return utils.RunCommand(command)
}

func (t *Tools) linkLibrary(ctx context.Context, output string) error {
	// https://blog.csdn.net/xuhongning/article/details/6365200
	args := []string{"-r", "-c"}
	if logs.IsDebug() {
		args = append(args, "-v")
	}
	args = append(args, t.LinkArgs...)
	args = append(args, output)
	for _, object := range t.objects {
		args = append(args, object)
	}
	command := exec.Command("ar", args...)
	command.Dir = t.Pwd
	logs.CtxDebug(ctx, "Link: %s", utils.PrettyCmd(command))
	return utils.RunCommand(command)
}

func (t *Tools) Run(ctx context.Context) error {
	runCmd := exec.Command(t.Output)
	runCmd.Dir = t.Pwd
	logs.CtxDebug(ctx, "Run: %s", utils.PrettyCmd(runCmd))
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
