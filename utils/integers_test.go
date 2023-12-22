package utils

import "testing"

func TestMaxInt(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "MaxInt(1, 2) = 2",
			args: args{x: 1, y: 2},
			want: 2,
		},
		{
			name: "MaxInt(2, 1) = 2",
			args: args{x: 2, y: 1},
			want: 2,
		},
		{
			name: "MaxInt(2, 2) = 2",
			args: args{x: 2, y: 2},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxInt(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("MaxInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinInt(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "MinInt(1, 2) = 1",
			args: args{x: 1, y: 2},
			want: 1,
		},
		{
			name: "MinInt(2, 1) = 1",
			args: args{x: 2, y: 1},
			want: 1,
		},
		{
			name: "MinInt(2, 2) = 2",
			args: args{x: 2, y: 2},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinInt(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("MinInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	type args struct {
		a string
		b int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Join(\"Hello\", 1) = \"Hello1\"",
			args: args{a: "Hello", b: 1},
			want: "Hello1",
		},
		{
			name: "Join(\"Hello\", 2) = \"Hello2\"",
			args: args{a: "Hello", b: 2},
			want: "Hello2",
		},
		{
			name: "Join(\"Hello\", 3) = \"Hello3\"",
			args: args{a: "Hello", b: 3},
			want: "Hello3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Join(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Join() = %v, want %v", got, tt.want)
			}
		})
	}
}
