package hoq

import (
	"HOQ/hoq/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_quicConnPool_add(t *testing.T) {
	assert := require.New(t)
	pool := &quicConnPool{}
	con1 := &quicConn{addr: ":7878", conn: &hmock.MockQuicSession{}}
	con2 := &quicConn{addr: ":8080", conn: &hmock.MockQuicSession{}}
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
