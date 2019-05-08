package test

import (
	"HOQ/hoq"
	"HOQ/logs"
	"testing"
)

/**
BenchmarkSpeed/QUIC-4         	    1000	   1924844 ns/op
BenchmarkSpeed/TCP-4          	    3000	    363133 ns/op
*/
func BenchmarkSpeed(b *testing.B) {
	qs, _ := hoq.NewServer(hoq.EngineQuic, hoq.EchoHandler)
	ts, _ := hoq.NewServer(hoq.EngineTcp, hoq.EchoHandler)
	go qs.Run("127.0.0.1:6665")
	go ts.Run("127.0.0.1:6667")
	qc, _ := hoq.NewClient(hoq.EngineQuic)
	tc, _ := hoq.NewClient(hoq.EngineTcp)
	logs.SetLevel(logs.LevelError)
	//quic 目前可能慢在TLS握手上，尝试去掉或给TCP加上TCP再次验证
	b.Run("QUIC", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := qc.Get("http://127.0.0.1:6665")
			if err != nil {
				logs.Error(err)
			}
		}
	})
	b.Run("TCP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := tc.Get("http://127.0.0.1:6667")
			if err != nil {
				logs.Error(err)
			}
		}
	})
}
