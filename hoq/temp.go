package hoq

/**
for temporary test&debug
this file may be deleted At any time of the
*/

const (
	testHost = "127.0.0.1:8080"
)

var testText = []byte{72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100, 33, 13}

var testRequestMessage = `POST http://detectportal.firefox.com/success.txt HTTP/1.1
Host: detectportal.firefox.com
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:66.0) Gecko/20100101 Firefox/66.0
Accept: */*
Accept-Language: zh-CN,en-US;q=0.7,en;q=0.3
Accept-Encoding: gzip, deflate
Cache-Control: no-cache
Pragma: no-cache
DNT: 1
Connection: keep-alive

this is body`
