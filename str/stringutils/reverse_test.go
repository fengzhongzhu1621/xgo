package stringutils

import "testing"

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
		})
	}
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
