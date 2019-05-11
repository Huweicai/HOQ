package ut

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommonPrefix(t *testing.T) {
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
			want: 2,
		},
		{
			args: args{
				a: "",
				b: "",
			},
			want: 0,
		},
		{
			args: args{
				a: "12345",
				b: "",
			},
			want: 0,
		},
		{
			args: args{
				a: "12345",
				b: "12366",
			},
			want: 3,
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

func TestDownSamplingDo(t *testing.T) {
	assert := require.New(t)
	x, y, z := 0, 0, 0
	total := 100000
	for i := 0; i < total; i++ {
		DownSamplingDo(0.1, func() {
			x++
		})
		DownSamplingDo(0, func() {
			y++
		})
		DownSamplingDo(1, func() {
			z++
		})
	}
	assert.True(x < 11000 && x > 9000)
	assert.Equal(total*0, y)
	assert.Equal(total*1, z)
}

func TestContain(t *testing.T) {
	target := "target"
	targets := []string{"hello", "world"}

	if Contain(nil, nil) || Contain(target, nil) || Contain(nil, targets) {
		t.Fail()
	}
	if Contain(target, targets) {
		t.Fail()
	}
	targets = nil
	if Contain(target, targets) {
		t.Fail()
	}
	targets = []string{"hello", "world", "target"}
	if !Contain(target, targets) {
		t.Fail()
	}

	mTargets := map[string]int{"hello": 0, "world": 0}
	if Contain(target, mTargets) {
		t.Fail()
	}

	mTargets = map[string]int{"hello": 0, "world": 0, "target": 0}
	if !Contain(target, mTargets) {
		t.Fail()
	}
}

func TestCommon(t *testing.T) {
	type args struct {
		targetsA []string
		targetsB []string
	}
	tests := []struct {
		name       string
		args       args
		wantResult []string
	}{
		{
			name: "test a", args: args{targetsA: []string{"x", "a", "b", "y"}, targetsB: []string{"t", "a", "b", "k"}}, wantResult: []string{"a", "b"},
		}, {
			name: "test b", args: args{targetsA: []string{"1", "2", "3", "4"}, targetsB: []string{"1", "9", "-1"}}, wantResult: []string{"1"},
		}, {
			name: "test c", args: args{targetsA: []string{"1", "2", "3", "4"}, targetsB: []string{}}, wantResult: nil,
		}, {
			name: "test d", args: args{targetsA: []string{"x", "a", "b", "y"}, targetsB: []string{"t"}}, wantResult: nil,
		}, {
			name: "test e", args: args{targetsA: []string{"x", "a", "b", "y"}, targetsB: []string{"a"}}, wantResult: []string{"a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := Common(tt.args.targetsA, tt.args.targetsB); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Common() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestSliceReduce(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name       string
		args       args
		wantDslice []string
	}{
		{
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"a"},
			},
			wantDslice: []string{"b", "c"},
		}, {
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"a", "b", "c"},
			},
			wantDslice: nil,
		}, {
			args: args{
				a: []string{},
				b: []string{"a", "b", "c"},
			},
			wantDslice: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDslice := SliceReduce(tt.args.a, tt.args.b); !reflect.DeepEqual(gotDslice, tt.wantDslice) {
				t.Errorf("SliceReduce() = %v, want %v", gotDslice, tt.wantDslice)
			}
		})
	}
}
