package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getRepoInfo(t *testing.T) {
	{
		info, err := getRepoInfo("git@github.com:golang/tools.git")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, info, &RepoInfo{Scheme: "ssh", Namespace: "github.com", Path: "golang/tools"})
		assert.Equal(t, info.Url(), "git@github.com:golang/tools.git")
	}
	{
		info, err := getRepoInfo("https://github.com/golang/tools.git")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, info, &RepoInfo{Scheme: "https", Namespace: "github.com", Path: "golang/tools"})
		assert.Equal(t, info.Url(), "https://github.com/golang/tools.git")
	}

	{
		info, err := getRepoInfo("http://github.com/golang/tools.git")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, info, &RepoInfo{Scheme: "http", Namespace: "github.com", Path: "golang/tools"})
		assert.Equal(t, info.Url(), "http://github.com/golang/tools.git")
	}

	{
		info, err := getRepoInfo("github.com/golang/tools")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, info, &RepoInfo{Scheme: "ssh", Namespace: "github.com", Path: "golang/tools"})
		assert.Equal(t, info.Url(), "git@github.com:golang/tools.git")
	}
}
