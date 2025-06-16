package hexo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/spf13/cobra"

	git "github.com/sabhiram/go-gitignore"

	"github.com/anthony-dong/golang/pkg/codec"
	"github.com/anthony-dong/golang/pkg/logs"
)

const (
	readmeFileName    = "README.md"
	gitIgnoreFileName = ".gitignore"
)

type markdownCommand struct {
	Dir          string         `json:"dir"`
	TemplateFile string         `json:"template"`
	Ignore       []string       `json:"git_ignore_pattern"`
	GitIgnore    *git.GitIgnore `json:"-"`
}

type readmeFileInfo struct {
	UrlPath string
	Total   int
}

func NewReadmeCmd() (*cobra.Command, error) {
	cmd := cobra.Command{Use: "readme", Short: "Generate a README file for a markdown project"}
	var (
		cfg = markdownCommand{}
	)
	cmd.Flags().StringVarP(&cfg.Dir, "dir", "d", "", "The markdown project dir")
	cmd.Flags().StringArrayVarP(&cfg.Ignore, "ignore", "i", nil, "The markdown template file path of gitignore pattern")
	cmd.Flags().StringVarP(&cfg.TemplateFile, "template", "t", "", fmt.Sprintf("The markdown template file path, go template struct: %+v", new(readmeFileInfo)))
	if err := cmd.MarkFlagRequired("dir"); err != nil {
		return nil, err
	}
	if err := cmd.MarkFlagRequired("template"); err != nil {
		return nil, err
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return cfg.Run(cmd.Context())
	}
	return &cmd, nil
}

func (m *markdownCommand) Run(ctx context.Context) error {
	if err := m.init(); err != nil {
		return err
	}
	parser, err := m.getParser()
	if err != nil {
		return err
	}
	logs.Info("New parser success, template file: %s", m.TemplateFile)

	info, err := m.getReadmeFileInfo()
	if err != nil {
		return err
	}
	logs.Info("Get ReadmeFileInfo success, Total: %d", info.Total)

	readmeFile := filepath.Join(m.Dir, readmeFileName)
	file, err := os.OpenFile(readmeFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	logs.Info("Open %s file success !!", readmeFile)

	if err := parser.Execute(file, info); err != nil {
		return err
	}
	logs.Info("Write README file success !!")
	return nil
}

func (m *markdownCommand) init() error {
	if filename, err := filepath.Abs(m.Dir); err != nil {
		return err
	} else {
		m.Dir = filename
	}
	if filename, err := filepath.Abs(m.TemplateFile); err != nil {
		return err
	} else {
		m.TemplateFile = filename
	}
	if ignoreFile := filepath.Join(m.Dir, gitIgnoreFileName); utils.ExistFile(filepath.Join(ignoreFile)) { // 不存在
		ignore, err := git.CompileIgnoreFileAndLines(ignoreFile, m.Ignore...)
		if err != nil {
			return err
		}
		m.GitIgnore = ignore
	} else {
		logs.Warn("[Markdown] not found %s in dir: %s", gitIgnoreFileName, m.Dir)
		m.GitIgnore = git.CompileIgnoreLines(m.Ignore...)
	}
	logs.Info("[Markdown] init config success:\n%s", utils.ToJson(m, true))
	return nil
}

func (m *markdownCommand) getParser() (*template.Template, error) {
	templateBody, err := os.ReadFile(m.TemplateFile)
	if err != nil {
		return nil, fmt.Errorf("read %s file find err: %v", m.TemplateFile, err)
	}
	temp := template.New("readme")
	parse, err := temp.Parse(string(templateBody))
	if err != nil {
		return nil, err
	}
	return parse, nil
}

func (m *markdownCommand) getReadmeFileInfo() (*readmeFileInfo, error) {
	builder := strings.Builder{}
	files, err := utils.GetAllFiles(m.Dir, func(fileName string) bool {
		if !(strings.HasSuffix(fileName, ".md") || strings.HasSuffix(fileName, ".markdown")) {
			return false
		}
		relativePath, err := filepath.Rel(m.Dir, fileName)
		if err != nil {
			logs.Warn("[Markdown] filepath.Rel(%s, %s) find err: %v", m.Dir, fileName, err)
			return false
		}
		if m.GitIgnore.MatchesPath(relativePath) {
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	// 转成 markdown的url写法
	toUrl := func(file string) string {
		file = strings.TrimPrefix(file, m.Dir)
		return fmt.Sprintf("- [%s](.%s)\n", file, string(codec.NewUrlCodec().Encode([]byte(file))))
	}
	//order := func(files []string) []string {
	//	l := skipset.NewString()
	//	for _, elem := range files {
	//		l.Add(elem)
	//	}
	//	result := make([]string, 0, len(files))
	//	l.Range(func(value string) bool {
	//		result = append(result, value)
	//		return true
	//	})
	//	return result
	//}
	// 遍历写
	for _, elem := range files {
		_, err := builder.WriteString(toUrl(elem))
		if err != nil {
			return nil, err
		}
	}
	return &readmeFileInfo{
		UrlPath: builder.String(),
		Total:   len(files),
	}, nil
}
