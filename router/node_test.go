package router

import (
	"HOQ/hoq"
	"HOQ/logs"
	"HOQ/util"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"testing"
)

var rt *node = NewSimpleNode("")

func init() {
	aNode := NewSimpleNode("hello/Germany")
	bNode := NewSimpleNode("hello/China")
	cNode := NewSimpleNode("hello/Canada")
	dNode := NewSimpleNode("hi/China")
	eNode := NewSimpleNode("hi/Canada")
	rt.addChildren(aNode, bNode, cNode, dNode, eNode,
		NewSimpleNode("hello/China/Hunan"),
		NewSimpleNode("hello/China/Hubei"),
		NewSimpleNode("hello/China/Hunan/Changsha"),
		NewSimpleNode("hello/China/Hunan/Changde"),
		NewSimpleNode("hello/Cuba"))
}

func TestNodes(t *testing.T) {
	assert := require.New(t)
	assert.NotPanics(func() {
		faNode := &node{
			val:      "",
			children: []*node{},
		}
		newNode := &node{
			val:      "helloW",
			children: []*node{},
		}
		newNode2 := &node{
			val:      "helloC",
			children: []*node{},
		}
		newNode3 := &node{
			val: "heXXX",
		}
		faNode.addChildren(nil)
		faNode.addChild(newNode)
		faNode.addChild(newNode2)
		faNode.addChild(newNode3)
	})
	fmt.Println(rt.print())
	fmt.Println("size", rt.size())

	nd := rt.find("hello/China/Hunan/Changsha")
	assert.NotNil(nd)
	p := nd.path()
	assert.Equal("hello/China/Hunan/Changsha", p)

	xNode := NewSimpleNode("hello")
	yNode := NewSimpleNode("nihao")
	xNode.addChild(yNode)
	//can't add node with out any relation except root ""
	assert.Equal(0, len(xNode.children))
	rt.sort()
	ut.Nothing()
	assert.True(rt.children[len(rt.children)-1].priority <= rt.children[0].priority)
}

func TestNodeMerge(t *testing.T) {
	a := NewSimpleNode("hello/G")
	a.methods = []string{"GET", "POST", "HEAD"}
	b := NewSimpleNode("hello/G")
	b.methods = []string{"GET", "POST", "DEL"}
	a.merge(b)
	assert.Equal(t, []string{"GET", "POST", "HEAD", "DEL"}, a.methods)

	r := NewSimpleNode(root)
	r.addChildren(a, b)
	assert.Equal(t, 1, len(r.children))
}

func TestNodeCopy(t *testing.T) {
	a := NewSimpleNode("hello/G")
	assert.Equal(t, a.val, a.copy().val)
}

func TestFind(t *testing.T) {
	got := rt.find("hello/Germany")
	logs.Info(got.val)
}

func TestList(t *testing.T) {
	assert := require.New(t)
	//空测试
	assert.NotPanics(func() {
		empty := nList([]*node{})
		empty.add(nil)
		assert.Nil(empty.get())
		assert.Nil(empty.peek())
	})

	aNode := NewSimpleNode("a")
	bNode := NewSimpleNode("b")
	nList := nList([]*node{})
	nList.add(aNode)
	assert.Equal(1, nList.len())
	nList.add(bNode)
	assert.Equal(2, nList.len())
	x := nList.peek()
	assert.Equal(2, nList.len())
	assert.Equal(len(nList), nList.len())
	nList.add(x)
	assert.Equal(3, nList.len())
	a := nList.get()
	assert.Equal(2, nList.len())
	assert.Equal("a", a.val)
	a = nList.get()
	assert.Equal(1, nList.len())
	assert.Equal("b", a.val)
	a = nList.get()
	assert.Equal(0, nList.len())
	assert.Equal("a", a.val)
}

/**
radix 树和哈希性能对比

BenchmarkNode/Radix-4         	 5000000	       358 ns/op
BenchmarkNode/Hash-4          	10000000	       151 ns/op
*/
func BenchmarkNode(b *testing.B) {
	testRadix := NewSimpleNode("")
	testHash := make(map[string]hoq.Handler)
	prefix := "hello/simple/prefix"
	total := 10000
	for i := 0; i < total; i++ {
		add(testRadix, testHash, prefix+strconv.Itoa(i))
	}
	b.Run("Radix", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := rand.Intn(total)
			s := prefix + strconv.Itoa(x)
			nd := testRadix.find(s)
			if nd == nil {
				logs.Error(s, "not found")
			}
		}
	})
	b.Run("Hash", func(b *testing.B) {
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
