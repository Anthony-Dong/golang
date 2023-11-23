package jsontool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/tidwall/gjson"

	"github.com/spf13/cobra"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "json",
		Short: "The Json tool",
		Long: `The Json Tool

  - Query JSON values: https://jqlang.github.io/jq
  - Terminal JSON viewer: https://github.com/antonmedv/fx
  - JSON Path: https://github.com/tidwall/gjson
`,
	}
	var (
		pretty bool
		path   string
		reader = os.Stdin
		writer = os.Stdout
	)
	cmd.Flags().StringVar(&path, "path", "", "set specified path")
	cmd.Flags().BoolVarP(&pretty, "pretty", "", false, "set pretty json")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if !utils.CheckStdInFromPiped() {
			return cmd.Help()
		}
		ctx := cmd.Context()
		out, err := readPath(ctx, reader, path)
		if err != nil {
			return err
		}
		if pretty {
			out, err = prettyJson(ctx, out)
			if err != nil {
				return err
			}
		}
		if _, err := writer.Write(out); err != nil {
			return err
		}
		return nil
	}
	return cmd, nil
}

func readPath(ctx context.Context, in io.Reader, path string) ([]byte, error) {
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	if path == "" {
		return data, nil
	}
	result := gjson.GetBytes(data, path)
	if !result.Exists() {
		return nil, fmt.Errorf(`not found path: %s`, path)
	}
	return []byte(result.String()), nil
}

func prettyJson(ctx context.Context, src []byte) ([]byte, error) {
	out := bytes.Buffer{}
	if err := json.Indent(&out, src, "", "  "); err != nil {
		return nil, err
	}
	if outBytes := out.Bytes(); len(outBytes) > 0 && outBytes[len(outBytes)-1] == '\n' {
		return outBytes, nil
	}
	out.WriteByte('\n')
	return out.Bytes(), nil
}
