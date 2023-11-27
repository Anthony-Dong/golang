package hexo

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"

	git "github.com/sabhiram/go-gitignore"
	"github.com/stretchr/testify/assert"

	"github.com/anthony-dong/golang/pkg/utils"
)

func TestGetAllPage(t *testing.T) {
	list, err := GetAllPage("test", []string{})
	if err != nil {
		t.Fatal(err)
	}
	for _, elem := range list {
		t.Log(elem)
	}
}
func TestLines(t *testing.T) {
	t.Run("绝对路径", func(t *testing.T) {
		lines := git.CompileIgnoreLines("/bin")
		assert.Equal(t, lines.MatchesPath("bin/tool"), true)
		assert.Equal(t, lines.MatchesPath("data/bin/tool"), false)
	})
	t.Run("相对路径", func(t *testing.T) {
		lines := git.CompileIgnoreLines(".git")
		assert.Equal(t, lines.MatchesPath("/.git/tool"), true)
		assert.Equal(t, lines.MatchesPath("/data/.git/tool"), true)
	})
}

func TestCheckFileCanHexoPre(t *testing.T) {
	assert.Equal(t, CheckFileCanHexoPre("test/hexo.md"), true)
	assert.Equal(t, CheckFileCanHexoPre("test/not_hexo.md"), false)
}

func TestCheckFileCanHexo(t *testing.T) {
	result, err := CheckFileCanHexo("test/hexo.md", "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(utils.ToJson(result, true))
}

func TestRun(t *testing.T) {
	dir := filepath.Clean("test")
	targetDir := filepath.Clean("test/post")
	firmCode := []string{"baidu", "ali"}
	if err := buildHexo(context.Background(), dir, targetDir, firmCode, nil); err != nil {
		t.Fatal(err)
	}
}

func TestYamlUnmarshal(t *testing.T) {
	data := "title: Maven构建工具介绍 \ndate: \"2022-11-15 20:06:04\"\ntags:\n\t- Java\ncategories:\n\t- Java\n\t- Maven"
	fmt.Println(data)
	data = `title: "Maven构建工具介绍" 
date: "2022-11-15 20:06:04"
tags:
  - "Java"
categories:
  - "Java"
  - "Maven"`
	config := Config{}
	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		t.Fatal(err)
	}
	t.Log(config)
}
