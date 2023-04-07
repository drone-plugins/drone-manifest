package tagging

import (
	"reflect"
	"testing"
)

func TestDefaultTags(t *testing.T) {
	tests := []struct {
		Before string
		After  []string
	}{
		// valid combinations
		{"", []string{"latest"}},
		{"refs/heads/master", []string{"latest"}},
		{"refs/tags/0.9.0", []string{"0.9", "0.9.0"}},
		{"refs/tags/1.0.0", []string{"1", "1.0", "1.0.0"}},
		{"refs/tags/v1.0.0", []string{"1", "1.0", "1.0.0"}},
		{"refs/tags/v1.0.0-alpha.1", []string{"1.0.0-alpha.1"}},

		// malformed or errors
		{"refs/tags/x1.0.0", []string{"latest"}},
		{"v1.0.0", []string{"latest"}},
		{"refs/tags/v18.06.0", []string{"18", "18.06", "18.06.0"}},
	}

	for _, test := range tests {
		got, want := DefaultTags(test.Before), test.After

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got tag %v, want %v", got, want)
		}
	}
}

func TestUseDefaultTag(t *testing.T) {
	type args struct {
		ref           string
		defaultBranch string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "latest tag for default branch",
			args: args{
				ref:           "refs/heads/master",
				defaultBranch: "master",
			},
			want: true,
		},
		{
			name: "build from tags",
			args: args{
				ref:           "refs/tags/v1.0.0",
				defaultBranch: "master",
			},
			want: true,
		},
		{
			name: "skip build for not default branch",
			args: args{
				ref:           "refs/heads/develop",
				defaultBranch: "master",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		if got := UseDefaultTag(tt.args.ref, tt.args.defaultBranch); got != tt.want {
			t.Errorf("%q. UseDefaultTag() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_stripHeadPrefix(t *testing.T) {
	type args struct {
		ref string
	}

	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				ref: "refs/heads/master",
			},
			want: "master",
		},
	}

	for _, tt := range tests {
		if got := stripHeadPrefix(tt.args.ref); got != tt.want {
			t.Errorf("stripHeadPrefix() = %v, want %v", got, tt.want)
		}
	}
}

func Test_stripTagPrefix(t *testing.T) {
	tests := []struct {
		Before string
		After  string
	}{
		{"refs/tags/1.0.0", "1.0.0"},
		{"refs/tags/v1.0.0", "1.0.0"},
		{"v1.0.0", "1.0.0"},
	}

	for _, test := range tests {
		got, want := stripTagPrefix(test.Before), test.After

		if got != want {
			t.Errorf("Got tag %s, want %s", got, want)
		}
	}
}
