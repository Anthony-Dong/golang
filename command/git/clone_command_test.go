package git

import (
	"reflect"
	"testing"
)

func Test_getRepoInfo(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    *RepoInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRepoInfo(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRepoInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRepoInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
