package hoq

import (
	"bou.ke/monkey"
	"crypto/tls"
	"errors"
	"github.com/lucas-clemente/quic-go"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
	"time"
)

func TestQuicEngine(t *testing.T) {
	assert := require.New(t)
	s, err := NewServer(EngineQuic, EchoHandler)
	assert.NoError(err)
	//是QUIC Engine
	assert.IsType(&quicEngine{}, s.engine)
	en, err := newEngine(EngineQuic, s)
	assert.IsType(&quicEngine{}, en)
	xErr := errors.New("test")
	guard := monkey.Patch(quic.ListenAddr,
		func(addr string, tlsConf *tls.Config, config *quic.Config) (quic.Listener, error) {
			return nil, xErr
		})
	assert.Equal(xErr, en.Serve(":8080"))
	guard.Unpatch()
	assert.NotPanics(func() {
		qc := en.(*quicEngine)
		go qc.Serve(":8080")
		time.Sleep(200 * time.Millisecond)
	})
}

func TestTcpEngine(t *testing.T) {
	assert := require.New(t)
	s, err := NewServer(EngineTcp, EchoHandler)
	assert.NoError(err)
	//tcp Engine
	assert.IsType(&tcpEngine{}, s.engine)
	en, err := newEngine(EngineTcp, s)
	assert.IsType(&tcpEngine{}, en)
	xErr := errors.New("test")
	guard := monkey.Patch(tls.Listen, func(network, laddr string, config *tls.Config) (net.Listener, error) {
		return nil, xErr
	})
	assert.Equal(xErr, en.Serve(":8080"))
	guard.Unpatch()
	assert.NotPanics(func() {
		go en.Serve(":8080")
		time.Sleep(200 * time.Millisecond)
	})
}
