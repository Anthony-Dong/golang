package run

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/anthony-dong/golang/command"

	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewCommand() (*cobra.Command, error) {
	filename := ""
	list := false
	debug := false
	vars := make([]string, 0)
	includes := make([]string, 0)
	run := &cobra.Command{
		Use:   "run",
		Short: `Run task templates`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := command.GetAppConfig(cmd.Context())
			if config.RunTaskConfig == nil {
				config.RunTaskConfig = &command.RunTaskConfig{}
			}
			helper := TaskRunner{}
			includes = append(includes, utils.GetPwd())
			if err := helper.InitIncludes(append(includes, config.RunTaskConfig.Includes...)); err != nil {
				return err
			}
			logs.Debug("includes: %s. enable debug: %v.", utils.ToJson(helper.Includes), debug)
			if list {
				for _, elem := range helper.Includes {
					if elem == "" {
						continue
					}
					files, err := utils.GetAllFilesWithMax(elem, func(fileName string) bool {
						isYaml := strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml")
						if !isYaml {
							return false
						}
						if _, err := helper.ParseTaskConfigFile(fileName); err != nil {
							return false
						}
						return true
					}, 64)
					if err != nil {
						return err
					}
					for _, file := range files {
						logs.Info("config file: %s", file)
					}
				}
				return nil
			}
			if filename == "" {
				return fmt.Errorf(`required flag(s) "config" not set`)
			}
			templateVars, err := helper.NewTemplateVars(vars)
			if err != nil {
				return err
			}
			logs.Info("init template vars: %s", utils.ToJson(templateVars, true))
			if filename, err = helper.LookupFile(filename, helper.Includes); err != nil {
				return err
			}
			configs, err := helper.ReadTaskConfig(filename, nil, templateVars)
			if err != nil {
				return err
			}
			if debug {
				logs.Info("file: %s\n%s", filename, utils.ToYaml(configs))
				return nil
			}
			logs.Info("file: %s\n%s", filename, utils.ToYaml(configs))
			for index, task := range configs {
				taskId := fmt.Sprintf("[task-%d]", index)
				run := exec.Command(task.Cmd, task.Args...)
				run.Dir = task.Dir
				run.Env = append(os.Environ(), task.Env...)
				logs.Info("%s start run %q", taskId, task.name())
				if err := utils.RunCmd(run, taskId+" ", task.Daemon); err != nil {
					return fmt.Errorf(`%s run %q retuen err: %v`, taskId, task.name(), err)
				}
				if task.Daemon {
					continue
				}
				logs.Info("%s end run %q", taskId, task.name())
			}
			return nil
		},
	}
	run.Flags().BoolVarP(&list, "list", "l", false, "List task config files")
	run.Flags().StringVarP(&filename, "config", "c", "", "The task config file")
	run.Flags().BoolVarP(&debug, "debug", "d", false, "Enable Debug")
	run.Flags().StringSliceVar(&vars, "var", []string{}, "Define template variables")
	run.Flags().StringSliceVarP(&includes, "include", "I", []string{}, "Add an config search path for includes")
	return run, nil
}

type TaskRunner struct {
	Includes []string
}

func (r *TaskRunner) InitIncludes(includes []string) error {
	skip := make(map[string]bool)
	includes = append(includes, utils.GetPwd())
	for _, elem := range includes {
		include, err := filepath.Abs(elem)
		if err != nil {
			return err
		}
		if skip[include] {
			continue
		}
		skip[include] = true
		r.Includes = append(r.Includes, include)
	}
	return nil
}

func (r *TaskRunner) NewTemplateVars(slice []string) (map[string]string, error) {
	kv := make(map[string]string, 0)
	for _, elem := range slice {
		kvs := strings.SplitN(elem, "=", 2)
		if len(kvs) != 2 {
			return nil, fmt.Errorf(`invalid template var: %q`, elem)
		}
		kv[strings.TrimSpace(kvs[0])] = strings.TrimSpace(kvs[1])
	}
	if kv["HOME"] == "" {
		kv["HOME"] = utils.GetUserHomeDir()
	}
	if kv["PWD"] == "" {
		kv["PWD"] = utils.GetPwd()
	}
	return kv, nil
}

type TaskConfig struct {
	Name    string            `json:"name,omitempty" yaml:"name,omitempty"`       // 任务名称
	Cmd     string            `json:"cmd,omitempty" yaml:"cmd,omitempty"`         // 命令名称，例如 ls/go/gcc
	Dir     string            `json:"dir,omitempty" yaml:"dir,omitempty"`         // 执行路径, 默认pwd
	Env     []string          `json:"env,omitempty" yaml:"env,omitempty"`         // 环境变量, 格式: k=v
	Args    []string          `json:"args,omitempty" json:"args,omitempty"`       // 命令参数
	Include string            `json:"include,omitempty" yaml:"include,omitempty"` // include 文件
	Run     string            `json:"run,omitempty" yaml:"run,omitempty"`         // 运行脚本，最终会替换成一个 shell命令
	Vars    map[string]string `json:"vars,omitempty" yaml:"vars,omitempty"`       // 模版变量
	Daemon  bool              `json:"daemon,omitempty" yaml:"daemon,omitempty"`   // 是否为后台任务
}

func (c *TaskConfig) name() string {
	if c.Name != "" {
		return c.Name
	}
	if c.Cmd != "" {
		return c.Cmd
	}
	if c.Run != "" {
		return c.Run
	}
	return "-"
}

func (c *TaskConfig) RenderTemplateVars(ctxValues map[string]string) (_ map[string]string, rErr error) {
	defer func() {
		if r := recover(); r != nil {
			if err := r.(error); err != nil {
				rErr = err
				return
			}
			rErr = fmt.Errorf(`%v`, r)
		}
	}()
	kv := utils.CopyMap(c.Vars)
	for k, v := range ctxValues {
		kv[k] = v
	}
	buffer := bytes.Buffer{}
	render := func(text string) string {
		if !strings.Contains(text, "{{") {
			return text
		}
		parse, err := template.New("").Option("missingkey=error").Funcs(map[string]interface{}{
			"IF": func(c1 string, c2 string, True string, False string) string {
				if c1 == c2 {
					return True
				}
				return False
			},
		}).Parse(text)
		if err != nil {
			panic(fmt.Errorf("parse temaplte:\n%s\nerr:\n%s", text, err.Error()))
		}
		buffer.Reset()
		if err := parse.Execute(&buffer, kv); err != nil {
			panic(fmt.Errorf("exec temaplte:\n%s\nerr:\n%s", text, err.Error()))
		}
		return buffer.String()
	}
	for key, value := range c.Vars {
		if vv, isExist := ctxValues[key]; isExist {
			c.Vars[key] = vv
			continue
		}
		kv[key] = render(value)
		c.Vars[key] = kv[key]
	}
	c.Name = render(c.Name)
	c.Cmd = render(c.Cmd)
	c.Dir = render(c.Dir)
	c.Include = render(c.Include)
	c.Run = render(c.Run)
	for index, elem := range c.Env {
		c.Env[index] = render(elem)
	}
	for index, elem := range c.Args {
		c.Args[index] = render(elem)
	}
	return kv, nil
}

func (*TaskRunner) ParseTaskConfig(content []byte) ([]*TaskConfig, error) {
	configs := make([]*TaskConfig, 0)
	if err := json.Unmarshal(content, &configs); err == nil {
		return configs, nil
	}
	decoder := yaml.NewDecoder(bytes.NewBuffer(content))
	for {
		childConfig := make([]*TaskConfig, 0)
		if err := decoder.Decode(&childConfig); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		configs = append(configs, childConfig...)
	}
	return configs, nil
}

func (r *TaskRunner) ParseTaskConfigFile(filename string) ([]*TaskConfig, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return r.ParseTaskConfig(file)
}

func (r *TaskRunner) IsEmpty(s string) bool {
	return s == "" || s == "-"
}

func (r *TaskRunner) IsNotEmpty(s string) bool {
	return !r.IsEmpty(s)
}

func (r *TaskRunner) ReadTaskConfig(filename string, walk map[string]bool, ctxValues map[string]string) ([]*TaskConfig, error) {
	abs, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf(`filepath.Abs(%q) find err: %v`, filename, err)
	}
	filename = abs
	configs, err := r.ParseTaskConfigFile(filename)
	if err != nil {
		return nil, err
	}
	if walk == nil {
		walk = map[string]bool{}
	}
	if walk[filename] {
		return nil, fmt.Errorf(`recycle %q file`, filename)
	}
	walk[filename] = true
	tasks := make([]*TaskConfig, 0)
	for _, task := range configs {
		taskVars, err := task.RenderTemplateVars(ctxValues)
		if err != nil {
			return nil, fmt.Errorf(`failed to render the template variables for task "%q", error: %v, config file: %s`, task.name(), err, filename)
		}
		if r.IsNotEmpty(task.Include) {
			includeFile, err := r.LookupIncludeFile(filename, task.Include, r.Includes)
			if err != nil {
				return nil, fmt.Errorf(`LookupIncludeFile(%q, %q) find err: %v`, filename, task.Include, err)
			}
			cmds, err := r.ReadTaskConfig(includeFile, walk, taskVars)
			if err != nil {
				return nil, err
			}
			tasks = append(tasks, cmds...)
			continue
		}
		for index, kv := range task.Env {
			key, value := utils.ReadKV(kv)
			if key == "" {
				return nil, fmt.Errorf(`invalid env name: %s, config file: %s`, kv, filename)
			}
			task.Env[index] = fmt.Sprintf("%s=%s", key, value) // 处理正确
		}
		if r.IsNotEmpty(task.Run) {
			tasks = append(tasks, &TaskConfig{
				Name: task.Name,
				Cmd:  "/bin/bash", // todo select env shell
				Args: []string{
					"-ce",
					task.Run,
				},
				Dir:  task.Dir,
				Env:  task.Env,
				Vars: task.Vars,
			})
			continue
		}
		if task.Cmd == "" {
			return nil, fmt.Errorf(`invald task command. config file: %s`, filename)
		}
		if filepath.Base(task.Cmd) == task.Cmd {
			if _, err := exec.LookPath(task.Cmd); err != nil {
				vErr, _ := err.(*exec.Error)
				if vErr != nil {
					return nil, fmt.Errorf(`exec.LookPath(%q) return err: %v, config file: %s`, task.Cmd, vErr.Err, filename)
				}
				return nil, fmt.Errorf(`exec.LookPath(%q) return err: %v, config file: %s`, task.Cmd, err, filename)
			}
		}
		args := make([]string, 0, len(task.Args))
		for _, arg := range task.Args {
			if arg == "-" {
				continue
			}
			args = append(args, arg)
		}
		task.Args = args
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// LookupIncludeFile
// case1:
// cur: /a/b/c.yaml
// filename: d.yaml
// return /a/b/d.yaml
func (r *TaskRunner) LookupIncludeFile(cur string, filename string, includes []string) (string, error) {
	if !filepath.IsAbs(cur) {
		abs, err := filepath.Abs(cur)
		if err != nil {
			return "", err
		}
		cur = abs
	}
	dir := filepath.Dir(cur)
	return r.LookupFile(filename, append([]string{dir}, includes...))
}

func (r *TaskRunner) LookupFile(filename string, includes []string) (string, error) {
	files := make([]string, 0)
	if filepath.IsAbs(filename) {
		files = append(files, filename)
	}
	for _, elem := range includes {
		files = append(files, filepath.Join(elem, filename))
	}
	//logs.Debug("files: %s, filename: %s", utils.ToJson(files), filename)
	for _, elem := range files {
		abs, err := filepath.Abs(elem)
		if err != nil {
			return "", err
		}
		if utils.ExistFile(abs) {
			return abs, nil
		}
	}
	return "", fmt.Errorf(`not found file: %s`, filename)
}
