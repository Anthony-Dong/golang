package golang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/utils"
)

func NewListGoModSize() (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "mod-size",
		Short: `list mod size`,
		RunE: func(cmd *cobra.Command, args []string) error {
			mod, err := ListMod(utils.GetPwd())
			if err != nil {
				return err
			}
			if err := SetModSize(mod); err != nil {
				return err
			}
			sort.Slice(mod, func(i, j int) bool {
				return mod[i].ZipSize < mod[j].ZipSize
			})
			totalSize := int64(0)
			for _, elem := range mod {
				totalSize = totalSize + elem.ZipSize
				fmt.Println(utils.ToJson(elem))
			}
			fmt.Println("totalSize:", totalSize)
			return nil
		},
	}, nil
}

type ModInfo struct {
	Path      string    `json:"Path"`
	Version   string    `json:"Version"`
	Time      time.Time `json:"Time"`
	Indirect  bool      `json:"Indirect"`
	GoMod     string    `json:"GoMod"`
	GoVersion string    `json:"GoVersion"`
	ZipSize   int64     `json:"ZipSize"`
}

func SetModSize(input []*ModInfo) error {
	for _, elem := range input {
		if elem.GoMod == "" {
			continue
		}
		filename := strings.TrimSuffix(elem.GoMod, filepath.Ext(elem.GoMod)) + ".zip"
		stat, err := os.Stat(filename)
		if err != nil {
			if strings.Contains(err.Error(), "no such file or directory") {
				continue
			}
			return fmt.Errorf(`state file [%s] find err: %s`, filename, err.Error())
		}
		elem.ZipSize = stat.Size()
	}
	return nil
}

func ListMod(dir string) ([]*ModInfo, error) {
	command := exec.Command("go", "list", "-m", "-json", "all")
	if dir != "" {
		command.Dir = dir
	}
	stdout := &bytes.Buffer{}
	errout := &bytes.Buffer{}
	command.Stdout = stdout
	command.Stderr = errout
	if err := command.Run(); err != nil {
		if s := stdout.String(); len(s) > 0 {
			fmt.Println(s)
		}
		fmt.Println(errout.String())
		return nil, err
	}
	decoder := json.NewDecoder(stdout)
	result := make([]*ModInfo, 0)
	for decoder.More() {
		m := &ModInfo{}
		if err := decoder.Decode(m); err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, nil
}
