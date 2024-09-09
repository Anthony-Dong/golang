package _init

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
)

type File struct {
	Name       string
	Content    string
	IsTemplate bool
}

func (f *File) Write(output string, replace bool) error {
	content, err := f.GetContent()
	if err != nil {
		return err
	}
	file := filepath.Join(output, f.Name)
	if utils.ExistFile(file) && !replace {
		logs.Info("file [%s] already exist", f.Name)
		return nil
	}
	logs.Info("create file [%s]", f.Name)
	if !utils.ExistDir(filepath.Dir(file)) {
		if err := os.MkdirAll(filepath.Dir(file), utils.DefaultDirMode); err != nil {
			return err
		}
	}
	if err := os.WriteFile(file, content, utils.DefaultFileMode); err != nil {
		return err
	}
	return nil
}

func (f *File) GetContent() ([]byte, error) {
	if !f.IsTemplate {
		return []byte(f.Content), nil
	}
	parse, err := template.New("").Parse(f.Content)
	if err != nil {
		return nil, err
	}
	output := bytes.Buffer{}
	if err := parse.Execute(&output, map[string]interface{}{}); err != nil {
		return nil, err
	}
	return nil, err
}
