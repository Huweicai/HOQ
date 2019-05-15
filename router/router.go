package router

import (
	"HOQ/hoq"
	"HOQ/logs"
	"HOQ/util"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	root       = ""
	ServerName = "HOQ-ROUTER"

	AdaptiveTrailingSlash = 1 << iota
)

var notFoundResp = hoq.FastResponse(hoq.StatusNotFound, defaultHeader())
var innerErrResp = hoq.FastResponse(hoq.StatusInternalServerError, defaultHeader())

/**
todo 限流及QPS统计
*/
type Router struct {
	adaptiveTrailingSlash bool
	addr                  string
	//为了性能，不加锁并发安全，一旦启动就进入只读模式，无法再修改
	started       bool
	root          *node
	bannedMethods []string
	//handler 执行前拦截器
	preInterFunc []PreInterceptor
	//handler 执行后拦截器
	sufInterFunc []SufInterceptor
	server       *hoq.Server
	//统计相关
	recordChan  chan *record
	records     map[string]int
	recordMutex sync.Mutex
}

/**
禁止某一种Method的使用
*/
func (t *Router) BanMethod(method string) bool {
	t.assertNotStarted()
	if !hoq.IsSupportedMethod(method) {
		return false
	}
	t.bannedMethods = append(t.bannedMethods, method)
	return true
}

func (t *Router) ShowRecordsCyclic(d time.Duration) {
	tk := time.Tick(d)
	for range tk {
		t.printRecords()
	}
}

/**
输出记录信息
*/
func (t *Router) printRecords() {
	t.recordMutex.Lock()
	defer t.recordMutex.Unlock()
	//methods
	methodRecords := make(map[string]int)
	pathRecords := make(map[string]int)
	respCodeRecords := make(map[string]int)
	for k, v := range t.records {
		if len(k) < 2 {
			logs.Error("unexpected short record key ", k)
			continue
		}
		switch key := k[2:]; k[:2] {
		case RecordPrefixMethod:
			methodRecords[key] = v
		case RecordPrefixPath:
			pathRecords[key] = v
		case RecordPrefixRespCode:
			respCodeRecords[key] = v
		default:
			logs.Error("unexpected record key ", k)
		}
	}
	logs.Info("method records\n" + m2String(methodRecords))
	logs.Info("url path records\n" + m2String(pathRecords))
	logs.Info("response code records\n" + m2String(respCodeRecords))
}

func m2String(m map[string]int) string {
	s := ""
	for k, v := range m {
		s += k + ":" + strconv.Itoa(v) + "\n"
	}
	return s
}

/**
记录信息，处理器
*/
func (t *Router) record() {
	for rcd := range t.recordChan {
		if rcd == nil {
			continue
		}
		t.recordMutex.Lock()
		t.records[rcd.k] += rcd.i
		t.recordMutex.Unlock()
	}
}

func (t *Router) AddPreInterceptor(itc PreInterceptor) {
	t.assertNotStarted()
	t.preInterFunc = append(t.preInterFunc, itc)
}

func (t *Router) AddSufInterceptor(itc SufInterceptor) {
	t.assertNotStarted()
	t.sufInterFunc = append(t.sufInterFunc, itc)
}

func (t *Router) assertNotStarted() {
	if t.started {
		panic("can not doing this after router started")
	}
}

/**
main router entrance
distribute requests to the handler they should be
*/
func (t *Router) main(ctx *hoq.Context) (rsp *hoq.Response) {
	//prevent the whole server form crashed by a panic of a request
	defer func() {
		if e := recover(); e != nil {
			logs.Error("something panicked ", e)
			rsp = innerErrResp
		}
	}()
	//前置拦截器
	for _, preFunc := range t.preInterFunc {
		err := preFunc(ctx, t)
		switch err := err.(type) {
		case nil:
			continue
		case *hoq.ErrWithCode:
			return hoq.FastResponse(err.Code(), nil)
		default:
			return innerErrResp
		}
	}
	//后置拦截器
	defer func() {
		for _, sufFunc := range t.sufInterFunc {
			err := sufFunc(rsp, t)
			switch err := err.(type) {
			case nil:
				continue
			case *hoq.ErrWithCode:
				rsp = hoq.FastResponse(err.Code(), nil)
			default:
				rsp = innerErrResp
			}
		}
	}()
	m := ctx.Request.Method()
	u := ctx.Request.URL()
	if ut.Contain(m, t.bannedMethods) {
		allows := ut.SliceReduce(hoq.Methods, t.bannedMethods)
		//banned
		return hoq.FastResponse(hoq.StatusMethodNotAllowed, defaultHeader().Set("Allow", strings.Join(allows, ",")))
	}
	handler := t.Find(m, u.Path)
	if handler == nil {
		return notFoundResp
	}
	return handler(ctx)
}

/**
新建一个HTTP router
*/
func New(httpNG hoq.NGType, addr string, flags int) (r *Router, err error) {
	r = &Router{root: NewSimpleNode(""), addr: addr}
	s, err := hoq.NewServer(httpNG, r.main)
	r.parseFlags(flags)
	if err != nil {
		return
	}
	r.server = s
	r.records = make(map[string]int)
	r.AddPreInterceptor(LogPreInterceptor)
	r.AddPreInterceptor(RecordPreInterceptor)
	r.AddSufInterceptor(RecordSufInterceptor)
	r.recordChan = make(chan *record, 1000)
	go r.record()
	return
}

/**
解析控制变量
*/
func (t *Router) parseFlags(flags int) {

}

/**
启动，开始执行，启动后就进入配置只读模式
*/
func (t *Router) Run() error {
	t.started = true
	t.root.sort()
	logs.Info("router started...")
	return t.server.Run(t.addr)
}

/*
Method GET wrapper for Add
*/
func (t *Router) GET(path string, handler hoq.Handler) error {
	return t.Add(path, handler, hoq.MethodGET)
}

/**
Method POST wrapper for Add
*/
func (t *Router) POST(path string, handler hoq.Handler) error {
	return t.Add(path, handler, hoq.MethodPOST)
}

/**
注册一个路径及对应Handler到Router
*/
func (t *Router) Add(path string, handler hoq.Handler, methods ...string) error {
	t.assertNotStarted()
	if handler == nil {
		return errors.New("handler cannot be empty")
	}
	for _, m := range methods {
		if !hoq.IsSupportedMethod(m) {
			return hoq.MethodNotSupportErr
		}
	}
	nd := &node{
		val:     path,
		methods: methods,
		handler: handler,
	}
	t.root.addChild(nd)
	return nil
}

func (t *Router) Find(method, path string) hoq.Handler {
	nd := t.root.find(path)
	if nd == nil {
		return nil
	}
	//判断方法是否支持
	if nd.allowAllMethods {
		return nd.handler
	}
	for _, m := range nd.methods {
		if m == method {
			return nd.handler
		}
	}
	logs.Info("found node ", nd.val)
	return nil
}

/**
新建一个Group
*/
func (t *Router) Group(path string) *Group {
	return &Group{
		prefix: path,
		t:      t,
	}
}
