package router

import (
	"HOQ/hoq"
	ut "HOQ/util"
	"strconv"
	"strings"
)

func NewSimpleNode(val string) *node {
	return &node{
		val: val,
	}
}

type node struct {
	father *node
	val    string
	//已注册支持的methods
	allowAllMethods bool
	methods         []string
	//处理器
	handler  hoq.Handler
	children []*node
}

func (n *node) hasSameSuffixWithChild(c *node) *node {
	for i, s := range n.children {
		if o := ut.CommonPrefix(c.val, s.val); o > 0 {
			return n.children[i]
		}
	}
	return nil
}

/**
hel
io llo  o

添加子节点

1 判断当前节点n和新节点c的value 关系
	a 完全无关
  		不应该让这种情况发生，父节点必须确认有关系才传递给子节点，
		对于根节点，"" 和 "abc" 视为有公共前缀""
	b n 和 c 有公共前缀
		b.1 n和c 有部分公共前缀
			分裂，为n创建新的字节点，n - common(n,c) , c - common(n,c)
		b.2 n包含于c
			为c查找包含剩余后缀的子节点
			b.2.1 找到
				递归
			b.2.2 找不到
				新建一个子节点

root.val = "" 所以万物都是root的子节点
*/
func (n *node) addChild(child *node) {
	if child == nil {
		return
	}
	i := ut.CommonPrefix(n.val, child.val)

	//a 没有无关，""为root
	if i == 0 && n.val != "" {
		return
	}
	//b
	//b.2 n 完全包含于child
	if i == len(n.val) {
		//真包含（重复的节点）
		if n.val == child.val {
			n.merge(child)
			return
		}
		child.cut(i)
		//b.2.1 找到了
		if e := n.hasSameSuffixWithChild(child); e != nil {
			e.addChild(child)
			return
		}
		//b.2.2 找不到
		n.linkChild(child)
		return
	}
	//b.1
	//原关键词分裂后的前缀，后缀
	pre := n.val[:i]
	suf := n.val[i:]
	//原有节点下沉一层，新增公共节点，并替换上层向下指针
	//child 截断
	child.cut(i)
	newDome := &node{
		father: n.father,
		val:    pre,
	}
	n.val = suf
	for i := range n.father.children {
		if n.father.children[i] == n {
			n.father.children[i] = newDome
		}
	}
	newDome.linkChild(child, n)
	return
}

func (n *node) getChildByVal(val string) *node {
	for _, c := range n.children {
		if c.val == val {
			return c
		}
	}
	return nil
}

func (n *node) addChildren(cs ...*node) {
	for _, c := range cs {
		n.addChild(c)
	}
}

/**
  just link whitout thinking instead of add (doing a lot of things)
*/
func (n *node) linkChild(c ...*node) {
	for _, t := range c {
		t.father = n
	}
	n.children = append(n.children, c...)
}

func (n *node) replaceChild(c ...*node) {
	for _, t := range c {
		t.father = n
	}
	n.children = c
}

/**
相同路径的节点合并（可能其他属性不同）
*/
func (n *node) merge(s *node) {

A:
	for _, smethod := range s.methods {
		for _, nmethod := range n.methods {
			if smethod == nmethod {
				continue A
			}
		}
		//没有在n中找到这个method
		n.methods = append(n.methods, smethod)
	}
}

/**
  从第i个位置开始截断value

  "abcd".cut(0) = "bcd"
*/
func (n *node) cut(i int) {
	if i < 0 {
		return
	}
	n.val = n.val[i:]
}

/**
  浅拷贝，父子指针值不变
*/
func (n *node) copy() *node {
	return &node{
		father:   n.father,
		val:      n.val,
		children: n.children,
	}
}

/**
  查找满足路径的子节点
  返回nil代表未找到
*/
func (n *node) findPrefixChild(path string) *node {
	remain := n.minusBy(path)
	for i, c := range n.children {
		if strings.HasPrefix(remain, c.val) {
			return n.children[i]
		}
	}
	return nil
}

/**
  返回请求路径减去当前节点路径的剩余路径
*/
func (n *node) minusBy(path string) string {
	if len(path) < len(n.val) {
		return path
	}
	return path[len(n.val):]
}

/**
  查找路径对应的节点：
  与当前节点比较
  A: 一样，当前节点即为目标节点
  B: 当前节点值完全包含于路径
  	a 找到是路径剩余部分的前缀子节点
  	b 找不到子节点：未能匹配
  C: 未能匹配
*/
func (n *node) find(path string) *node {
	//A
	if n.val == path {
		return n
	}
	i := ut.CommonPrefix(n.val, path)
	//B
	if i == len(n.val) {
		c := n.findPrefixChild(path)
		//B.b
		if c == nil {
			return nil
		}
		//B.a
		return c.find(n.minusBy(path))
	}
	//C
	return nil
}

/**
返回从头至尾的value路径
*/
func (n *node) path() string {
	if n.father == nil {
		return ""
	}
	return n.father.path() + n.val
}

/**
总结点个数
*/
func (n *node) size() (total int) {
	//self
	total += 1
	for _, c := range n.children {
		total += c.size()
	}
	return
}

/**
层次遍历，输出为string
*/
func (n *node) print() (out string) {
	list := nList([]*node{})
	list.add(n)
	num := 0
	for {
		l := list.len()
		if l == 0 {
			break
		}
		out += strconv.Itoa(num)
		num++
		for i := 0; i < l; i++ {
			e := list.get()
			out += " " + e.val + " "
			for x := range e.children {
				list.add(e.children[x])
			}
		}
		out += "\n"
	}
	return
}

/*
自定义实现的node list
*/
type nList []*node

func (l *nList) add(n *node) {
	if n == nil {
		return
	}
	*l = append(*l, n)
}

func (l *nList) get() (n *node) {
	if len(*l) == 0 {
		return nil
	}
	n = (*l)[0]
	*l = (*l)[1:]
	return
}

func (l *nList) peek() (n *node) {
	if len(*l) == 0 {
		return nil
	}
	return (*l)[0]
}

func (l *nList) len() int {
	return len(*l)
}
