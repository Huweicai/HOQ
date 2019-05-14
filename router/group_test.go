package router

import (
	"HOQ/hoq"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGroup(t *testing.T) {
	assert := require.New(t)
	r, err := New(hoq.EngineQuic, "127.0.0.1:8787", 0)
	assert.NoError(err)
	r.Add("/hello", hoq.EchoHandler, hoq.MethodGET)
	g := r.Group("/find")
	g.Add("/666", hoq.EchoHandler, hoq.MethodGET)
	g.GET("/777", hoq.EchoHandler)
	g.POST("/888", hoq.EchoHandler)
	assert.NotNil(r.Find(hoq.MethodGET, "/find/666"))
	assert.NotNil(r.Find(hoq.MethodGET, "/find/777"))
	assert.NotNil(r.Find(hoq.MethodPOST, "/find/888"))
}
