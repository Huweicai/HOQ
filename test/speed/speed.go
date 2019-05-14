package speed

import (
	"HOQ/hoq"
	"HOQ/logs"
	"sync"
	"time"
)

const testTxt = "sfdjsakfasdkjfsdjalfsklfkadslfadsfdsafjsdlfadsjkfklncvmxcn,vafadlfjqweofwejfldskfjksadflkasdflkdsfjsdjflskdjfdskfsjalfjdsfjaklsdfaskljdsafjdakfjdsafjsdkflafjksadjfladsnvcmvnewfjwejsdkflksjsdklfdsfnvc,vnsdafjdlsfjkdsfjksdlfjsadkfadsjf"

func main() {
	qs, _ := hoq.NewServer(hoq.EngineQuic, hoq.EchoHandler)
	ts, _ := hoq.NewServer(hoq.EngineTcp, hoq.EchoHandler)
	go qs.Run("10.8.125.150:6665")
	go ts.Run("10.8.125.150:6667")
	qc, _ := hoq.NewClient(hoq.EngineQuic)
	tc, _ := hoq.NewClient(hoq.EngineTcp)
	logs.SetLevel(logs.LevelWarn)
	//quic 目前可能慢在TLS握手上，尝试去掉或给TCP加上TCP再次验证
	wg := sync.WaitGroup{}
	wg.Add(2)
	var n = 100
	logs.Warn("======测试开始======")
	go func() {
		defer wg.Done()
		start := time.Now()
		for i := 0; i < n; i++ {
			_, err := qc.Post("http://10.8.125.150:6665", []byte(testTxt))
			if err != nil {
				logs.Error(err)
			}
		}
		end := time.Now().Sub(start).Nanoseconds()
		logs.Warn("QUIC", n, "次，平均耗时", end, "ns")
	}()
	go func() {
		defer wg.Done()
		start := time.Now()
		for i := 0; i < n; i++ {
			_, err := tc.Post("http://10.8.125.150:6667", []byte(testTxt))
			if err != nil {
				logs.Error(err)
			}
		}
		end := time.Now().Sub(start).Nanoseconds()
		logs.Warn("TCP ", n, "次，平均耗时", end, "ns")
	}()
	wg.Wait()
}
