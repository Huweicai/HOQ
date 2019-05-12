package hoq

import (
	"HOQ/logs"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_quicConnPool_add(t *testing.T) {
	assert := require.New(t)
	pool := &quicConnPool{}
	con1 := &quicConn{addr: ":7878", conn: &MockQuicSession{}}
	con2 := &quicConn{addr: ":8080", conn: &MockQuicSession{}}
	pool.add(con1)
	pool.add(con2)
	assert.Equal(2, pool.size())
	g1, err := pool.get(con1.addr)
	assert.NoError(err)
	assert.Equal(con1.addr, g1.addr)
	g2, err := pool.get(con2.addr)
	assert.NoError(err)
	assert.Equal(con2.addr, g2.addr)
}

func TestMockDemo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessClient := NewMockQuicSession(ctrl)
	//注册一个Mock方法
	sessClient.EXPECT().Close().Return(errors.New("error failed"))
	//正常调用
	err := sessClient.Close()
	logs.Error(err)
	assert.Error(t, err)
	assert.Equal(t, "error failed", err.Error())
}
