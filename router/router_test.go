package router

import (
	"HOQ/hoq"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRouter(t *testing.T) {
	assert := require.New(t)
	//empty test
	assert.NotPanics(func() {
		r, err := New(hoq.EngineQuic, "127.0.0.1:8786", 0)
		assert.NoError(err)
		assert.Error(r.Add("", nil))
		assert.Error(r.Add("", hoq.EchoHandler, "666"))
	})
	r, err := New(hoq.EngineQuic, "127.0.0.1:8787", 0)
	assert.NoError(err)
	r.Add("/hello", hoq.EchoHandler, hoq.MethodGET)
	assert.NotNil(r.Find(hoq.MethodGET, "/hello"))
	assert.NotPanics(func() {
		go r.Run()
	})
}
