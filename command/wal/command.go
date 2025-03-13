package wal

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/wal"
)

func NewCommand() (*cobra.Command, error) {
	var filename string
	var key string
	cmd := &cobra.Command{
		Use:   "wal",
		Short: "the wal command",
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := wal.OpenFile(filename)
			if err != nil {
				return err
			}
			w, err := wal.NewWal(file, 1024*1024)
			if err != nil {
				return err
			}
			if key != "" {
				get, err := w.Get(key)
				if err != nil {
					return err
				}
				fmt.Println(string(get))
				return nil
			}
			for _, elem := range w.List() {
				fmt.Println(elem)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&filename, "file", "f", "", "file name")
	cmd.Flags().StringVarP(&key, "key", "k", "", "file name")
	if err := cmd.MarkFlagRequired("file"); err != nil {
		return nil, err
	}
	return cmd, nil
}
