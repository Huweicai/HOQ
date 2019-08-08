package router

import (
	"HOQ/hoq"
	"HOQ/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
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
	go r.Run()
	assert.NotNil(r.Find(hoq.MethodGET, "/find/666"))
	assert.NotNil(r.Find(hoq.MethodGET, "/find/777"))
	assert.NotNil(r.Find(hoq.MethodPOST, "/find/888"))
	c, _ := hoq.NewClient(hoq.EngineQuic)
	c.SetReqTimeout(5 * time.Minute)
	ctx, err := c.Get("http://127.0.0.1:8787/find/888")
	ut.Nothing(ctx, err)

}
