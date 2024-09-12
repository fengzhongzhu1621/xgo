package xgo

import (
	"reflect"
	"sort"
	"testing"
)

func TestVersionList(t *testing.T) {
	v1 := VersionList{
		{1, 2, 3},
		{1, 2, 4},
		{1, 1, 3},
		{2, 1, 3},
	}
	sort.Sort(v1)
	expected := VersionList{
		{1, 1, 3},
		{1, 2, 3},
		{1, 2, 4},
		{2, 1, 3},
	}
	if !reflect.DeepEqual(v1, expected) {
		t.Errorf("VersionList sorting failed. Expected %v, got %v", expected, v1)
	}
}

func TestGetLatestVersion(t *testing.T) {
	type args struct {
		filelist []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "正确匹配",
			args: args{filelist: []string{"1.2.3", "1.2.4", "1.1.3", "2.1.3"}},
			want: "V2.1.3",
		},
		{
			name: "错误匹配",
			args: args{filelist: []string{"bad version"}},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLatestVersion(tt.args.filelist); got != tt.want {
				t.Errorf("getLatestVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
