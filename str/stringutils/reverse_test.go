package stringutils

import (
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
)

func TestReverseString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "testing empty",
			args: args{
				s: "",
			},
			want: "",
		},
		{
			name: "testing not empty",
			args: args{
				s: "hello world",
			},
			want: "dlrow olleh",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseString(tt.args.s); got != tt.want {
				t.Errorf("ReverseString() = %v, want %v", got, tt.want)
			}
			if got := strutil.Reverse(tt.args.s); got != tt.want {
				t.Errorf("strutil.Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}
