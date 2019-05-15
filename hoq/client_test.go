package hoq

import (
	"HOQ/logs"
	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
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

func TestClient(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c, err := NewClient(EngineQuic)
	c.SetReqTimeout(1 * time.Second)
	guard := monkey.Patch((*QUICCourier).RoundTrip, func(_ *QUICCourier, req *Request) (*Response, *RemoteInfo, error) {
		logs.Info("begin sleep")
		time.Sleep(3 * time.Second)
		logs.Info("sleep finished")
		return nil, nil, nil
	})
	defer guard.Unpatch()
	_, err = c.Get("http://127.0.0.1:8080")
	assert.Equal(RequestTimeoutErr, err)
}
