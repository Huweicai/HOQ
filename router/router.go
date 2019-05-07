package router

import "HOQ/util"

type RadixTree struct {
}

type node struct {
	father   *node
	children []*node
	val      string
}

func (n *node) hasSameSuffixWithChild(c *node) *node {
	for i, s := range n.children {
		if o := ut.CommonPrefix(c.val, s.val); o >= 0 {
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
	//a 没有无关
	if i == 0 {
		return
	}
	//b
	//b.2 n 完全包含于child
	if i == len(n.val) {
		//b.2.1 找到了
		if e := n.hasSameSuffixWithChild(child); e != nil {
			e.addChild(child)
			return
		}
		//b.2.2 找不到
		n.children = append(n.children, child)
		return
	}
	//b.1
	//分裂后的前缀
	pre := n.val[:i+1]
	//分裂后原来的节点
	suf := n.val[i:]
	newN := n.copy()
	newN.val = suf
	newN.father = n
	n.val = pre
	n.children = []*node{newN, child}
	return
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
