package hoq

import (
	"github.com/stretchr/testify/require"
	"net/textproto"
	"testing"
)

func Test_readGotCookies(t *testing.T) {
	assert := require.New(t)
	h := textproto.MIMEHeader(map[string][]string{
		"Cookies": {"key1=123;key2=234;key3=456"},
	})
	got := readGotCookies(h)
	assert.Equal(3, len(got))
	assert.Equal("key1", got[0].Name)
	assert.Equal("123", got[0].Value)
}

func Test_readSetCookies(t *testing.T) {
	assert := require.New(t)
	h := textproto.MIMEHeader(map[string][]string{
		"Set-Cookie": {"cookie-3=three; Domain=example.com", "cookie-2=two; Max-Age=3600"},
	})
	got := readSetCookies(h)
	assert.Equal(2, len(got))
	assert.Equal("cookie-3", got[0].Name)
	assert.Equal("three", got[0].Value)
}

var writeSetCookiesTests = []struct {
	Cookie *Cookie
	Raw    string
}{
	{
		&Cookie{Name: "cookie-1", Value: "v$1"},
		"cookie-1=v$1",
	},
	{
		&Cookie{Name: "cookie-2", Value: "two", MaxAge: 3600},
		"cookie-2=two; Max-Age=3600",
	},
	{
		&Cookie{Name: "cookie-3", Value: "three", Domain: ".example.com"},
		"cookie-3=three; Domain=example.com",
	},
}

func TestCookie_String(t *testing.T) {
	for _, cookie := range writeSetCookiesTests {
		if cookie.Raw != cookie.Cookie.String() {
			t.Fail()
		}
	}
}
