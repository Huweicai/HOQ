package hoq

import (
	"HOQ/logs"
	"io"
	"net/url"
	"strings"
	"testing"
)

func Test_parseFirstLine(t *testing.T) {
	a := []byte("\n")
	logs.Error(a)
	type args struct {
		line string
	}
	tests := []struct {
		name       string
		args       args
		wantMethod string
		wantUrl    string
		wantProto  string
		wantOk     bool
	}{
		{args: args{""},
			wantMethod: "",
			wantUrl:    "",
			wantProto:  "",
			wantOk:     false,
		},
		{args: args{"GET /foo HTTP/1.1"},
			wantMethod: "GET",
			wantUrl:    "/foo",
			wantProto:  "HTTP/1.1",
			wantOk:     true,
		},
		{args: args{"POST http://example.com HTTP/2.0"},
			wantMethod: "POST",
			wantUrl:    "http://example.com",
			wantProto:  "HTTP/2.0",
			wantOk:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMethod, gotUrl, gotProto, gotOk := parseFirstRequestLine(tt.args.line)
			if gotMethod != tt.wantMethod {
				t.Errorf("parseFirstLine() gotMethod = %v, want %v", gotMethod, tt.wantMethod)
			}
			if gotUrl != tt.wantUrl {
				t.Errorf("parseFirstLine() gotUrl = %v, want %v", gotUrl, tt.wantUrl)
			}
			if gotProto != tt.wantProto {
				t.Errorf("parseFirstLine() gotProto = %v, want %v", gotProto, tt.wantProto)
			}
			if gotOk != tt.wantOk {
				t.Errorf("parseFirstLine() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

/*
BenchmarkParseFirstLine/bad-4         	20000000	        86.5 ns/op
BenchmarkParseFirstLine/good-4        	100000000	        17.5 ns/op
4 times faster than the normal one
*/
func BenchmarkParseFirstLine(b *testing.B) {
	b.Run("bad", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			parseFirstLineBad("POST http://example.com HTTP/2.0")
		}
	})
	b.Run("good", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			parseFirstRequestLine("POST http://example.com HTTP/2.0")
		}
	})

}

func parseFirstLineBad(line string) (method, url, proto string, ok bool) {
	got := strings.Split(line, " ")
	if len(got) == 3 {
		return got[0], got[1], got[2], true
	}
	return
}

func TestRequest_requestLine(t *testing.T) {
	type fields struct {
		method  string
		url     *url.URL
		proto   string
		headers *Headers
		Body    io.Reader
	}
	u1, _ := url.Parse("http://www.example.com")
	u2, _ := url.Parse("/Home")
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{},
			want:   "  ",
		},
		{
			fields: fields{
				proto:  "HTTP/1.1",
				url:    u1,
				method: "GET",
			},
			want: "GET http://www.example.com HTTP/1.1",
		},
		{
			fields: fields{
				proto:  "HTTP/2.0",
				url:    u2,
				method: "POST",
			},
			want: "POST /Home HTTP/2.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				method:  tt.fields.method,
				url:     tt.fields.url,
				proto:   tt.fields.proto,
				headers: tt.fields.headers,
				Body:    tt.fields.Body,
			}
			if got := r.requestLine(); got != tt.want {
				t.Errorf("Request.requestLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
