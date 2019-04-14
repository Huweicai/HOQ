package hoq

import (
	"github.com/lucas-clemente/quic-go"
	"net"
)

/**
一条连接，在TCP中是一条连接，而在QUIC中是一个Stream吗？
*/
type Channel interface {
	net.Conn
}

type quicChannel struct {
	core quic.Session
}

func newQuicChannel(session quic.Session) *quicChannel {
	return &quicChannel{core: session}
}
