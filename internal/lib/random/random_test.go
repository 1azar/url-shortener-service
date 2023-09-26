package random

import "testing"

func TestNewRandomString(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test", args{size: 6}, "xxxxx"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRandomString(tt.args.size); got != tt.want {
				t.Errorf("NewRandomString() = %v, want %v", got, tt.want)
			}
		})
	}
}
