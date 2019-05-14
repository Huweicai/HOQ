package router

import (
	"HOQ/hoq"
	"HOQ/logs"
	"HOQ/util"
	"math/rand"
	"strconv"
	"testing"
)

/**
radix 树和哈希性能对比
*/
func BenchmarkNode(b *testing.B) {
	testRadix := NewSimpleNode("")
	testHash := make(map[string]hoq.Handler)
	prefix := "hello/simple/prefix"
	total := 10000
	for i := 0; i < total; i++ {
		add(testRadix, testHash, prefix+strconv.Itoa(i))
	}
	b.Run("Hash", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := rand.Intn(total)
			s := prefix + strconv.Itoa(x)
			nd := testRadix.find(s)
			if nd == nil {
				logs.Error(s, "not found")
			}
		}
	})
	b.Run("Radix", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := rand.Intn(total)
			s := prefix + strconv.Itoa(x)
			got, ok := testHash[s]
			if ok != true {
				logs.Error(s, "not found")
			}
			ut.Nothing(got)
		}
	})
}

func add(n *node, m map[string]hoq.Handler, s string) {
	n.addChild(NewSimpleNode(s))
	m[s] = nil
}
