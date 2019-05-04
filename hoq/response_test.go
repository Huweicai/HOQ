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
