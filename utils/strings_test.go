package utils

import "testing"

func TestTruncate(t *testing.T) {
	type args struct {
		s string
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Truncate long string",
			args: args{s: "Hello World", n: 10},
			want: "Hello[...]",
		},
		{
			name: "Do not truncate short string",
			args: args{s: "Hello World", n: 20},
			want: "Hello World",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Truncate(tt.args.s, tt.args.n); got != tt.want {
				t.Errorf("Truncate() = %v, want %v", got, tt.want)
			}
		})
	}
}
