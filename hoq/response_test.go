package hoq

import (
	"io"
	"testing"
)

func TestResponse_statusLine(t *testing.T) {
	type fields struct {
		proto      string
		statusCode int
		statusMSg  string
		headers    *Headers
		Body       io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				proto:      "HTTP/1.1",
				statusMSg:  "OK",
				statusCode: 200,
			},
			want: "HTTP/1.1 200 OK",
		},
		{
			want: "  ",
		},
		{
			fields: fields{
				proto:      "HTTP/2.0",
				statusCode: 400,
				statusMSg:  "Bad Request",
			},
			want: "HTTP/2.0 400 Bad Request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				proto:      tt.fields.proto,
				statusCode: tt.fields.statusCode,
				statusMSg:  tt.fields.statusMSg,
				headers:    tt.fields.headers,
				Body:       tt.fields.Body,
			}
			if got := r.statusLine(); got != tt.want {
				t.Errorf("Response.statusLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseStatusLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name      string
		args      args
		wantCode  int
		wantMsg   string
		wantProto string
		wantOk    bool
	}{
		{
			args: args{
				"HTTP/2.0 400 Bad Request",
			},
			wantCode:  400,
			wantMsg:   "Bad Request",
			wantProto: "HTTP/2.0",
			wantOk:    true,
		},
		{
			args: args{
				"HTTP/1.1 200 OK",
			},
			wantCode:  200,
			wantMsg:   "OK",
			wantProto: "HTTP/1.1",
			wantOk:    true,
		},
		{
			args:   args{},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode, gotMsg, gotProto, gotOk := parseStatusLine(tt.args.line)
			if gotCode != tt.wantCode {
				t.Errorf("parseStatusLine() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
			if gotMsg != tt.wantMsg {
				t.Errorf("parseStatusLine() gotMsg = %v, want %v", gotMsg, tt.wantMsg)
			}
			if gotProto != tt.wantProto {
				t.Errorf("parseStatusLine() gotProto = %v, want %v", gotProto, tt.wantProto)
			}
			if gotOk != tt.wantOk {
				t.Errorf("parseStatusLine() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
