package hoq

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func Test_bodyLength(t *testing.T) {
	type args struct {
		r io.Reader
	}
	bys := []byte("Hello World")
	tests := []struct {
		name       string
		args       args
		wantLength int64
		wantOk     bool
	}{
		{
			args: args{
				bytes.NewReader(bys),
			},
			wantLength: int64(len(bys)),
			wantOk:     true,
		}, {
			args: args{
				nil,
			},
			wantLength: 0,
			wantOk:     true,
		}, {
			args: args{
				NoBody,
			},
			wantLength: 0,
			wantOk:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLength, gotOk := bodyLength(tt.args.r)
			if gotLength != tt.wantLength {
				t.Errorf("bodyLength() gotLength = %v, want %v", gotLength, tt.wantLength)
			}
			if gotOk != tt.wantOk {
				t.Errorf("bodyLength() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_bodyAllowedForStatus(t *testing.T) {
	type args struct {
		status int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				100,
			},
			want: false,
		},
		{
			args: args{
				204,
			},
			want: false,
		},
		{
			args: args{
				304,
			},
			want: false,
		},
		{
			args: args{
				200,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bodyAllowedForStatus(tt.args.status); got != tt.want {
				t.Errorf("bodyAllowedForStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_urlParse(t *testing.T) {
	assert := require.New(t)
	assert.NotPanics(func() {
		u, err := urlParse("http://example.com")
		assert.NoError(err)
		assert.Equal("example.com", u.Host)
		assert.Equal("http", u.Scheme)
		_, err = urlParse("ftp://example.com")
		assert.Error(err)
	})
}
