package main

import (
	"HOQ/hoq"
	"HOQ/logs"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

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
	var n, txtSize = 30, 9000
	logs.Warn("======测试开始======")
	go func() {
		defer wg.Done()
		start := time.Now()
		fail := 0
		for i := 0; i < n; i++ {
			var testTxt = RandStringRunes(txtSize)
			_, err := qc.Post("http://10.8.125.150:6665", []byte(testTxt))
			if err != nil {
				fail++
			}
		}
		cost := time.Now().Sub(start).Nanoseconds()
		logs.Warn(fmt.Sprintf("QUIC %d次请求 请求载荷大小：%dB 平均耗时：%dns 失败率:%f%%", 1000, txtSize, cost/int64(n), float32(0)/float32(n)))
	}()
	go func() {
		defer wg.Done()
		start := time.Now()
		fail := 0
		for i := 0; i < n; i++ {
			var testTxt = RandStringRunes(9000)
			_, err := tc.Post("http://10.8.125.150:6667", []byte(testTxt))
			if err != nil {
				fail++
			}
		}
		cost := time.Now().Sub(start).Nanoseconds()
		logs.Warn(fmt.Sprintf("TCP  %d次请求 请求载荷大小：%dB 平均耗时：%dns 失败率:%f%%", 1000, txtSize, cost/int64(n), float32(0)/float32(n)))
	}()
	wg.Wait()
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
