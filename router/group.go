package router

import "HOQ/hoq"

/**
一组方法，拥有共同的URL前缀
*/
type Group struct {
	prefix string
	t      *Router
}

func (g *Group) Add(path string, handler hoq.Handler, methods ...string) error {
	path = g.prefix + path
	return g.t.Add(path, handler, methods...)
}

/*
Method GET wrapper for Add
*/
func (t *Group) GET(path string, handler hoq.Handler) error {
	return t.Add(path, handler, hoq.MethodGET)
}

/**
Method POST wrapper for Add
*/
func (t *Group) POST(path string, handler hoq.Handler) error {
	return t.Add(path, handler, hoq.MethodPOST)
}
