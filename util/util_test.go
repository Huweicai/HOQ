package ut

import "testing"

func TestCommon(t *testing.T) {
	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				a: "abc",
				b: "ab",
			},
			want: 1,
		},
		{
			args: args{
				a: "",
				b: "",
			},
			want: -1,
		},
		{
			args: args{
				a: "12345",
				b: "",
			},
			want: -1,
		},
		{
			args: args{
				a: "12345",
				b: "12366",
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CommonPrefix(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Common() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				a: 1123123,
				b: 2999999,
			},
			want: 1123123,
		}, {
			args: args{},
			want: 0,
		}, {
			args: args{
				a: 100,
				b: 100,
			},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMax(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				a: 1123123,
				b: 2999999,
			},
			want: 2999999,
		}, {
			args: args{},
			want: 0,
		}, {
			args: args{
				a: 100,
				b: 100,
			},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}
