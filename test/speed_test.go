package test

import (
	"HOQ/hoq"
	"HOQ/logs"
	"testing"
)

const testTxt = "sfdjsakfasdkjfsdjalfsklfkadslfadsfdsafjsdlfadsjkfklncvmxcn,vafadlfjqweofwejfldskfjksadflkasdflkdsfjsdjflskdjfdskfsjalfjdsfjaklsdfaskljdsafjdakfjdsafjsdkflafjksadjfladsnvcmvnewfjwejsdkflksjsdklfdsfnvc,vnsdafjdlsfjkdsfjksdlfjsadkfadsjf"

/**
第一次测试，quic tls VS 纯tcp
BenchmarkSpeed/QUIC-4         	    1000	   1924844 ns/op
BenchmarkSpeed/TCP-4          	    3000	    363133 ns/op

第二次测试 quic + tls VS tcp + tls
BenchmarkSpeed/QUIC-4         	     500	   2005767 ns/op
BenchmarkSpeed/TCP-4          	     500	   2451111 ns/op

第三次测试 同步quic tcp tls 配置
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
			_, err := qc.Post("http://127.0.0.1:6665", []byte(testTxt))
			if err != nil {
				logs.Error(err)
			}
		}
	})
	b.Run("TCP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := tc.Post("http://127.0.0.1:6667", []byte(testTxt))
			if err != nil {
				logs.Error(err)
			}
		}
	})
}
