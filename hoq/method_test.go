package hoq

import "testing"

func Test_isSupportedMethod(t *testing.T) {
	type args struct {
		method string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{MethodGET}, want: true,
		},
		{
			args: args{"TEST"}, want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSupportedMethod(tt.args.method); got != tt.want {
				t.Errorf("isSupportedMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
