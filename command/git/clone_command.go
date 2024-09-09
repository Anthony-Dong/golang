package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
)

func NewGitCloneCommand() (*cobra.Command, error) {
	home := ""
	cmd := &cobra.Command{
		Use:   "clone [url] [--branch branch] [--depth 1]",
		Short: "fast git clone repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf(`invalid repo name`)
			}
			gitInfo, err := getRepoInfo(args[0])
			if err != nil {
				return err
			}
			cloneArgs := []string{"clone"}
			if len(args) > 1 {
				cloneArgs = append(cloneArgs, args[1:]...)
			}
			homeAbs, err := filepath.Abs(home)
			if err != nil {
				return fmt.Errorf(`filepath.Abs(%q) find err: %v`, home, err)
			}
			cloneArgs = append(cloneArgs, gitInfo.CloneUrl(), gitInfo.CloneDir(homeAbs))
			execCmd := exec.Command("git", cloneArgs...)
			logs.StdOut(utils.PrettyCmd(execCmd))
			return utils.RunCommand(execCmd)
		},
	}
	cmd.Flags().StringVar(&home, "home", filepath.Join(utils.GetUserHomeDir(), "go/src"), "The home dir")
	return cmd, nil
}

type RepoInfo struct {
	Scheme    string // ssh/https
	Namespace string
	Path      string
}

func (r *RepoInfo) CloneUrl() string {
	switch r.Scheme {
	case "ssh":
		return fmt.Sprintf("git@%s:%s.git", r.Namespace, r.Path)
	default:
		return fmt.Sprintf("%s://%s/%s.git", r.Scheme, r.Namespace, r.Path)
	}
}

func (r *RepoInfo) CloneDir(home string) string {
	return filepath.Join(home, strings.ToLower(r.Namespace), strings.ToLower(filepath.Dir(r.Path)), filepath.Base(r.Path))
}

func getRepoInfo(url string) (*RepoInfo, error) {
	// https://github.com/golang/tools.git
	// git@github.com:golang/tools.git
	// github.com/golang/tools
	if strings.HasPrefix(url, "git@") {
		if !regexp.MustCompile(`^git@[a-zA-Z0-9._-]+:[a-zA-Z0-9/._-]+\.git$`).MatchString(url) {
			return nil, fmt.Errorf(`invalid git url: %s`, url)
		}
		split := strings.Split(strings.TrimSuffix(strings.TrimPrefix(url, "git@"), ".git"), ":")
		if len(split) != 2 {
			return nil, fmt.Errorf(`invalid git url: %s`, url)
		}
		return &RepoInfo{Scheme: "ssh", Namespace: split[0], Path: split[1]}, nil
	}
	if strings.HasPrefix(url, "https://") {
		url = strings.TrimSuffix(strings.TrimPrefix(url, "https://"), ".git")
		split := strings.Split(url, "/")
		if len(split) != 3 {
			return nil, fmt.Errorf(`invalid git url: %s`, url)
		}
		return &RepoInfo{Scheme: "https", Namespace: split[0], Path: split[1] + "/" + split[2]}, nil
	}
	if strings.HasPrefix(url, "http://") {
		url = strings.TrimSuffix(strings.TrimPrefix(url, "http://"), ".git")
		split := strings.Split(url, "/")
		if len(split) != 3 {
			return nil, fmt.Errorf(`invalid git url: %s`, url)
		}
		return &RepoInfo{Scheme: "http", Namespace: split[0], Path: split[1] + "/" + split[2]}, nil
	}
	split := strings.Split(url, "/")
	if len(split) != 3 {
		return nil, fmt.Errorf(`invalid git url: %s`, url)
	}
	return &RepoInfo{Scheme: "ssh", Namespace: split[0], Path: split[1] + "/" + split[2]}, nil
}
