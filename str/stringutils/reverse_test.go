package stringutils

import (
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gookit/goutil/arrutil"
	"github.com/stretchr/testify/assert"
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

func TestReflectReverseSlice(t *testing.T) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	ReflectReverseSlice(names)
	expected := []string{"g", "f", "e", "d", "c", "b", "a"}
	assert.Equal(t, expected, names)
}

func TestReverseSliceGetNew(t *testing.T) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	newNames := ReverseSliceGetNew(names)
	expected := []string{"g", "f", "e", "d", "c", "b", "a"}
	assert.Equal(t, expected, newNames)
}

func TestReverseSlice(t *testing.T) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	ReverseSlice(names)
	expected := []string{"g", "f", "e", "d", "c", "b", "a"}
	assert.Equal(t, expected, names)
}

func BenchmarkReverseReflectSlice(b *testing.B) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	for i := 0; i < b.N; i++ {
		ReflectReverseSlice(names)
	}
}

func BenchmarkReverseSlice(b *testing.B) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	for i := 0; i < b.N; i++ {
		ReverseSlice(names)
	}
}

func BenchmarkReverseSliceNew(b *testing.B) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	for i := 0; i < b.N; i++ {
		ReverseSliceGetNew(names)
	}
}

func TestReverse(t *testing.T) {
	ss := []string{"a", "b", "c"}
	arrutil.Reverse(ss)
	assert.Equal(t, []string{"c", "b", "a"}, ss)

	ints := []int{1, 2, 3}
	arrutil.Reverse(ints)
	assert.Equal(t, []int{3, 2, 1}, ints)
}
