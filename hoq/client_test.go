package hoq

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewClient(t *testing.T) {
	assert := require.New(t)
	c, err := NewClient(EngineQuic)
	assert.NoError(err)
	assert.IsType(&QUICCourier{}, c.engine)
	c, err = NewClient(EngineTcp)
	assert.NoError(err)
	assert.IsType(&TCPCourier{}, c.engine)
}
