package random

import "testing"

func TestNewRandomString(t *testing.T) {
	type args struct {
		size uint
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test", args{size: 6}, "xxxxxx"},
		{"test", args{size: 1}, "x"},
		{"test", args{size: 3}, "xxx"},
		{"test", args{size: 9}, "xxxxxxxxx"},
		{"test", args{size: 11}, "xxxxxxxxxxx"},
		{"test", args{size: 0}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRandomString(tt.args.size); len(got) != len(tt.want) {
				t.Errorf("NewRandomString() = %v, want %v", got, tt.want)
			}
		})
	}
}
