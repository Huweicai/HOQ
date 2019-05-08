package router

import (
	"HOQ/hoq"
	"HOQ/logs"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestRouter(t *testing.T) {
	assert := require.New(t)
	//empty test
	assert.NotPanics(func() {
		r, err := New(hoq.EngineQuic, "127.0.0.1:8786", 0)
		assert.NoError(err)
		assert.Error(r.Add("", nil))
		assert.Error(r.Add("", hoq.EchoHandler, "666"))
	})
	r, err := New(hoq.EngineQuic, "127.0.0.1:8787", 0)
	assert.NoError(err)
	r.Add("/hello", hoq.EchoHandler, hoq.MethodGET)
	assert.NotNil(r.Find(hoq.MethodGET, "/hello"))
	assert.NotPanics(func() {
		go r.Run()
	})
}
