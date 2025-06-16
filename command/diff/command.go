package diff

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/anthony-dong/golang/command"
	"github.com/anthony-dong/golang/pkg/diff"
	"github.com/anthony-dong/golang/pkg/utils"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Diff utilities for various data formats",
	}
	if err := command.AddCommand(cmd, NewJsonDiff); err != nil {
		return nil, err
	}
	return cmd, nil
}

func NewJsonDiff() (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "json",
		Short: "Diff JSON files or data",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 3 {
				file, err := os.ReadFile(args[0])
				if err != nil {
					return err
				}
				return diffJsonPath(file, args[1], args[2])
			}
			if len(args) == 2 {
				return diffJson(args[0], args[1])
			}
			return fmt.Errorf(`unsupport args %#v`, args)
		},
	}, nil
}

func diffJson(file1, file2 string) error {
	json1, err := os.ReadFile(file1)
	if err != nil {
		return err
	}
	json2, err := os.ReadFile(file2)
	if err != nil {
		return err
	}
	diffs, err := diff.DiffJson(json1, json2)
	if err != nil {
		return err
	}
	for _, elem := range diffs {
		key := elem.Key
		if !strings.HasPrefix(key, ".") {
			key = "." + key
		}
		fmt.Printf("%q\n", key)
		fmt.Println(utils.ToJson(elem.Origin))
		fmt.Println(utils.ToJson(elem.Patch))
		fmt.Println("------------------------------------------------------------")
	}
	return nil
}

func diffJsonPath(jsonData []byte, jsonPath1 string, jsonPath2 string) error {
	json1 := gjson.GetBytes(jsonData, jqPath2GJson(jsonPath1))
	if json1.Type == gjson.Null {
		return fmt.Errorf(`"%s" is null`, jsonPath1)
	}
	json2 := gjson.GetBytes(jsonData, jqPath2GJson(jsonPath2))
	if json2.Type == gjson.Null {
		return fmt.Errorf(`"%s" is null`, jsonPath2)
	}
	diffs, err := diff.DiffJsonString(json1.Raw, json2.Raw)
	if err != nil {
		return err
	}
	for _, elem := range diffs {
		key := ""
		switch elem.Type {
		case diff.DelType, diff.ChangeType:
			key = jsonPath1 + elem.Key
		case diff.AddType:
			key = jsonPath2 + elem.Key
		}
		if !strings.HasPrefix(key, ".") {
			key = "." + key
		}
		fmt.Printf("%q\n", key)
	}
	return nil
}

func gJsonPath2Jq(path string) string {
	compile := regexp.MustCompile(`\.(\d+)\.`)
	return compile.ReplaceAllString(path, ".[${1}].")
}

func jqPath2GJson(path string) string {
	compile := regexp.MustCompile(`\.\[(\d+)]\.`)
	return strings.TrimPrefix(compile.ReplaceAllString(path, ".${1}."), ".")
}
