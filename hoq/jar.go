package hoq

import (
	"net/url"
	"strings"
	"sync"
	"time"
)

type CookieJar interface {
	// SetCookies handles the receipt of the cookies in a reply for the
	// given URL.  It may or may not choose to save the cookies, depending
	// on the jar's policy and implementation.
	SetCookies(u *url.URL, cookies []*Cookie)

	// Cookies returns the cookies to send in a request for the given URL.
	// It is up to the implementation to honor the standard cookie use
	// restrictions such as in RFC 6265.
	GetCookies(u *url.URL) []*Cookie
}

type cookieJar struct {
	lock sync.Mutex
	//核心存储结构：一级域名，子域名，path
	m map[string]map[string][]*Cookie
}

var jar *cookieJar

func init() {
	jar = &cookieJar{
		m: make(map[string]map[string][]*Cookie),
	}
	go jar.removeExpires()
}

/**
周期性的清理过期Cookie
*/
func (c *cookieJar) removeExpires() {
	tk := time.NewTicker(10 * time.Second)
	for range tk.C {
		c.lock.Lock()
		for _, topdM := range c.m {
			for d, dM := range topdM {
				var newCks []*Cookie
				for _, ck := range dM {
					if ck.expired() {
						continue
					}
					newCks = append(newCks, ck)
				}
				topdM[d] = newCks
			}
		}
		c.lock.Unlock()
	}
}

func newCookieJava() {
	return
}

func (c *cookieJar) SetCookies(u *url.URL, cookies []*Cookie) {
	topd := topDomain(u.Host)
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.m[topd] == nil {
		n := make(map[string][]*Cookie)
		n[u.Host] = cookies
	} else if len(c.m[topd][u.Host]) == 0 {
		c.m[topd][u.Host] = cookies
	} else {
		c.m[topd][u.Host] = append(c.m[topd][u.Host], cookies...)
	}
}

/**
从CookieJar中根据指定的URL取出对应的Cookie
*/
func (c *cookieJar) GetCookies(u *url.URL) (cookies []*Cookie) {
	topd := topDomain(u.Host)
	path := u.Path
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.m[topd] == nil {
		return
	}
	if len(c.m[topd][u.Host]) == 0 {
		return
	}
	cks := c.m[topd][u.Host]
	for _, ck := range cks {
		if ck.expired() {
			continue
		}
		if !ck.matchPath(path) {
			continue
		}
		cookies = append(cookies, ck)
	}
	return
}

/**
返回域名的顶级域名
www.example.com -> example.com
*/
func topDomain(d string) string {
	i := strings.LastIndex(d, ".")
	j := strings.LastIndex(d[:i], ".")
	return d[j+1:]
}
