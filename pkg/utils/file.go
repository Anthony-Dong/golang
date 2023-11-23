package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultFileMode os.FileMode = 0644
	DefaultDirMode  os.FileMode = 0755
	FileSeparator               = filepath.Separator
)

// GetGoProjectDir 有 go.mod 的目录.
func GetGoProjectDir() string {
	path := filepath.Dir(os.Args[0])
	if ExistFile(filepath.Join(path, "go.mod")) {
		return path
	}
	path, err := os.Getwd()
	if err == nil {
		max := 4
		cur := 0
		for cur < max {
			if ExistFile(filepath.Join(path, "go.mod")) {
				return path
			}
			path = filepath.Dir(path)
			cur++
		}
	}
	return "."
}

// Exist 判断文件是否存在.
func ExistFile(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func GetFilePrefixAndSuffix(filename string) (prefix, suffix string) {
	filename = filepath.Base(filename)
	ext := filepath.Ext(filename)
	if ext == "" {
		return filename, ""
	}
	filename = strings.TrimSuffix(filename, ext)
	return filename, ext
}

// foo 不能写操作，引用需要deepcopy
func ReadLineByFunc(file io.Reader, foo func(line string) error) error {
	if file == nil {
		return fmt.Errorf("ReadLines find reader is nil")
	}
	reader := bufio.NewReader(file)
	prefix := make([]byte, 0)
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if isPrefix {
			prefix = append(prefix, line...)
			continue
		}
		if len(prefix) > 0 {
			prefix = append(prefix, line...)
			line = prefix
			prefix = prefix[:0]
		}
		if err := foo(string(line)); err != nil {
			return err
		}
	}
	return nil
}

func ReadLines(read io.Reader) ([]string, error) {
	result := make([]string, 0)
	if err := ReadLineByFunc(read, func(line string) error {
		result = append(result, line)
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

type FilterFile func(fileName string) bool

// GetAllFiles 从路径dirPth下获取全部的文件.
func GetAllFiles(dirPth string, filter FilterFile) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(dirPth, func(path string, info os.FileInfo, err error) error {
		if info != nil && info.IsDir() {
			return nil
		}
		if filter(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// GetFileRelativePath fileName指的是文件的路径 path 指的是文件的父路径地址，return 相对路径.
func GetFileRelativePath(fileName string, path string) (string, error) {
	//return filepath.Rel(path, fileName)
	var err error
	if fileName, err = filepath.Abs(fileName); err != nil {
		return "", err
	}
	if path, err = filepath.Abs(path); err != nil {
		return "", err
	}
	// 没有前缀说明不在目录
	if !strings.HasPrefix(fileName, path) {
		return "", fmt.Errorf("the file %v not in path %v", fileName, path)
	}
	relativePath := strings.TrimPrefix(fileName, path)
	relativePath = filepath.Clean(relativePath)
	if strings.HasPrefix(relativePath, string(filepath.Separator)) {
		return filepath.Clean(strings.TrimPrefix(relativePath, string(filepath.Separator))), nil
	}
	return relativePath, nil
}

func WriteFile(filename string, content []byte) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil
	}
	defer file.Close()
	if content == nil {
		content = []byte{}
	}
	if _, err := file.Write(content); err != nil {
		return err
	}
	return nil
}

func GetCmdName() string {
	//return "go-tool"
	return strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
}

func MustTmpDir(dir string, pattern string) string {
	if dir, err := ioutil.TempDir(dir, pattern); err != nil {
		panic(err)
	} else {
		return dir
	}
}

func UserHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "/root"
	}
	return dir
}

// CheckStdInFromPiped 检测标准输入是否来自于管道符
func CheckStdInFromPiped() bool {
	if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
		return true
	} else {
		return false
	}
}
