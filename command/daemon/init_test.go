package bstatic

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/anthony-dong/golang/pkg/utils"
)

type Process struct {
	Name string
	Bin  string
	Args []string
	Desc string
	Envs map[string]string
}

func TestName(t *testing.T) {
	file, err := ReadFile(LinuxDaemonFile)
	if err != nil {
		t.Error(err)
	}
	template := utils.MustTemplate("test", nil, file)

	out := bytes.Buffer{}
	err = template.Execute(&out, Process{
		Bin:  "go",
		Name: "my_test.service",
		Args: []string{"run", "main.go"},
		Desc: "test process",
		Envs: map[string]string{
			"TEST_ENV_1": "1",
			"TEST_ENV_2": "2",
		},
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(out.String())

	// 写入到 /etc/systemd/system/ 目录下
	//sudo systemctl list-units --type=service # 列出全部的service
	//sudo systemctl enable my_test.service    # 启动并开启服务
	//sudo systemctl start my_test.service     # 启动服务
	//sudo systemctl restart my_test.service   # 重启服务
	//sudo systemctl status my_test.service    # 检查服务现状
	//sudo journalctl -u my_test.service       # 查看日志
}
