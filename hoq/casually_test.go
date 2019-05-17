package hoq

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMsg(t *testing.T) {
	assert := require.New(t)
	req, _ := NewRequest(MethodGET,
		"http://127.0.0.1:8080",
		NewHeaders(map[string]string{"Want": "X"}), NoBody)
	resp, _ := NewResponse(StatusOK, nil, NoBody)
	testMsg := func(msg Message) {
		msg.SetHeader(testHeader0)
		assert.Equal(testHeader0, msg.GetHeader())
		msg.SetBody(NoBody)
		assert.Equal(NoBody, msg.GetBody())
	}
	testMsg(resp)
	testMsg(req)
	assert.NoError(req.EatFirstLine("GET /foo HTTP/1.1"))
	assert.Equal("GET", req.Method())
	assert.Equal("GET /foo HTTP/1.1", req.FirstLine())

	assert.NoError(resp.EatFirstLine("HTTP/1.1 200 OK"))
	assert.Equal(StatusOK, resp.Code())
	assert.Equal("HTTP/1.1 200 OK", resp.FirstLine())
}

func TestTls(t *testing.T) {
	tls := generateQuicTLSConfig()
	assert.NotNil(t, tls)
}

func TestHandler(t *testing.T) {
	assert := require.New(t)
	req, _ := NewRequest(MethodGET,
		"http://127.0.0.1:8080",
		NewHeaders(map[string]string{"Want": "X"}), NoBody)
	ctx := &Context{
		Request: req,
	}
	resp := EchoHandler(ctx)
	assert.Equal(StatusOK, resp.statusCode)

	resp = ByeHandler(ctx)
	assert.Equal(StatusOK, resp.statusCode)
}

const urlQuic = "127.0.0.1:7002"
const urlTcp = "127.0.0.1:7001"

func TestAll(t *testing.T) {
	go startTwoServer()
	time.Sleep(100 * time.Millisecond)
	qc, _ := NewClient(EngineQuic)
	tc, _ := NewClient(EngineTcp)
	testClientSend(qc, urlQuic)
	testClientSend(tc, urlTcp)
	time.Sleep(100 * time.Millisecond)
}
func testClientSend(c *Client, u string) {
	c.SetReqTimeout(100 * time.Millisecond)
	ctx, err := c.Get("http://" + u)
	ctx, err = c.Post("http://"+u, []byte("hello world"))
	if err != nil {
		panic(err)
	}
	if ctx.Response.Code() != 200 {
		panic("not 200")
	}
}

func startTwoServer() {
	qserver, err := NewServer(EngineQuic, EchoHandler)
	if err != nil {
		panic(err)
	}
	tserver, err := NewServer(EngineTcp, EchoHandler)
	if err != nil {
		panic(err)
	}
	go qserver.Run(urlQuic)
	go tserver.Run(urlTcp)
}

func TestStatusLine(t *testing.T) {
	assert.Equal(t, "HTTP/1.1 200 OK\r\n", string(statusLine(StatusOK)))
}
