package hoq

import (
	"HOQ/util"
	"github.com/lucas-clemente/quic-go"
	"sync"
	"time"
)

type quicConn struct {
	addr  string
	conn  quic.Session
	cTime *time.Time
	uTime *time.Time
}

type QuicSession interface {
	quic.Session
}

type quicConnPool struct {
	lock sync.Mutex
	pool map[string][]*quicConn
}

/**
尝试从缓存中获取一条连接，如果没有的话，就创建一条
*/
func (p *quicConnPool) get(addr string) (conn *quicConn, err error) {
	p.lock.Lock()
	got, ok := p.pool[addr]
	if ok && len(got) != 0 {
		conn = got[0]
		conn.uTime = ut.Now()
		//peek the first one ,add to the last
		p.pool[addr] = got[1:]
		p.lock.Unlock()
		return
	}
	p.lock.Unlock()
	return p.make(addr)
}

/**
return a quic connection
*/
func (p *quicConnPool) add(conn *quicConn) bool {
	//lazy init
	if p.pool == nil {
		p.pool = make(map[string][]*quicConn)
	}
	if conn == nil || conn.conn == nil || conn.addr == "" {
		return false
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.pool[conn.addr] = append(p.pool[conn.addr], conn)
	return true
}

func (p *quicConnPool) make(addr string) (conn *quicConn, err error) {
	sess, err := quic.DialAddr(addr, generateTCPTLSConfig(), nil)
	return &quicConn{addr: addr, conn: sess, uTime: ut.Now(), cTime: ut.Now()}, err
}

func (p *quicConnPool) size() int {
	return len(p.pool)
}
