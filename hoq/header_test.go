package hoq

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/textproto"
	"strings"
	"testing"
)

var testHeader0 = NewHeaders(nil)
var testHeader1 = NewHeaders(map[string]string{"Host": "example.com"})
var testHeader2 = NewHeaders(map[string]string{"Host": "example.com", "Content-Length": "10"})
var testHeader3 = NewHeaders(map[string]string{"Host": "example.com", "Content-Length": "9", "User-Agent": "Firefox 57", "": ""})

func TestHeaders_GenContentLength(t *testing.T) {
	assert := require.New(t)
	var t2 = NewHeaders(map[string]string{"Host": "example.com", "Content-Length": "10"})
	l1 := t2.ContentLength()
	assert.Equal(int64(10), l1)
	b1 := []byte{1, 2, 3, 4, 5, 6}
	ok := t2.GenContentLength(bytes.NewReader(b1))
	assert.True(ok)
	l2 := t2.ContentLength()
	assert.Equal(int64(6), l2)

	b2 := "123456789"
	ok = t2.GenContentLength(strings.NewReader(b2))
	assert.True(ok)
	l3 := t2.ContentLength()
	assert.Equal(int64(9), l3)
}

func TestHeaders_GetSet(t *testing.T) {
	assert := require.New(t)
	var h = NewHeaders(map[string]string{"Host": "example.com", "Content-Length": "10"})
	assert.Equal("", h.Get("aaa"))
	h.Set("aaa", "bbb")
	assert.Equal("bbb", h.Get("aaa"))
	assert.True(h.Exits("aaa"))
}

func TestHeaders_Serialize(t *testing.T) {
	got := testHeader3.Serialize()
	assert.Equal(t, len("Host: example.com\r\nContent-Length: 9\r\nUser-Agent: Firefox 57"), len(got))
}

func TestMimeToMap(t *testing.T) {
	assert.NotPanics(t, func() {
		mimeHeaderToMap(nil)
	})
	mm := map[string][]string{
		"1": []string{"1", "2", "3"},
		"2": []string{"1", "2", "3"},
		"3": []string{"1", "2", "3"},
		"4": nil,
	}

	m := mimeHeaderToMap(textproto.MIMEHeader(mm))
	assert.Equal(t, 3, len(m))
	assert.Equal(t, "1", m["1"])
	assert.Equal(t, "1", m["2"])
	assert.Equal(t, "1", m["3"])
}

func TestConvertHTTPHeader(t *testing.T) {
	assert.NotPanics(t, func() {
		convertHttpHeader(nil)
	})
	mm := map[string][]string{
		"1": []string{"1", "2", "3"},
		"2": []string{"1", "2", "3"},
		"3": []string{"1", "2", "3"},
		"4": nil,
	}
	m := convertHttpHeader(http.Header(mm))
	assert.Equal(t, 3, len(m.headers))
	assert.Equal(t, "1", m.Get("1"))
	assert.Equal(t, "1", m.Get("2"))
	assert.Equal(t, "1", m.Get("3"))
}
