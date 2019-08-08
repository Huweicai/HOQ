package hoq

import (
	"testing"
)

func Test_topDomain(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				d: "www.example.com",
			},
			want: "example.com",
		},
		{
			args: args{
				d: "example.com",
			},
			want: "example.com",
		}, {
			args: args{
				d: "example.www.example.com",
			},
			want: "example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := topDomain(tt.args.d); got != tt.want {
				t.Errorf("topDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
