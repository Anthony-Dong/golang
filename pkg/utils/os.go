package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var (
	gopath     string
	gopathOnce sync.Once
)

const (
	gopathName = "GOPATH"
)

func GetGoPath() string {
	gopathOnce.Do(func() {
		// load  env
		gopath = os.Getenv(gopathName)
		if gopath != "" {
			return
		}

		// load go env
		stdOut := bytes.Buffer{}
		command := exec.Command("go", "env", gopathName)
		command.Stdout = &stdOut
		if err := command.Run(); err == nil {
			gopath = strings.TrimSuffix(stdOut.String(), "\n")
		}

		// load default home
		if gopath == "" {
			dir, err := os.UserHomeDir()
			if err == nil {
				gopath = filepath.Join(dir, "go")
			}
		}

		// load default name
		if gopath == "" {
			gopath = "/go"
		}
	})

	return gopath
}

var _GetUserHomeDirOnce sync.Once
var _UserHomeDir = ""

func GetUserHomeDir() string {
	_GetUserHomeDirOnce.Do(func() {
		dir, err := os.UserHomeDir()
		if err != nil {
			panic(fmt.Errorf(`os.UserHomeDir return err: %v`, err))
		}
		_UserHomeDir = dir
	})
	return _UserHomeDir
}

var _GetPwd string
var _GetPwdOnce sync.Once

func GetPwd() string {
	_GetPwdOnce.Do(func() {
		dir, err := os.Getwd()
		if err != nil {
			panic(fmt.Errorf(`os.Getwd return err: %v`, err))
		}
		_GetPwd = dir
	})
	return _GetPwd
}
