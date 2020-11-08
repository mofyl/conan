package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"conan/config"
	"conan/log"
)

const (
	MAX_MEM = 32 << 20 // 32M
	CLOSED  = 1
	START   = 2
)

var (
	default404Body = []byte("404 page not Found")
	default405Body = []byte("405 method not allowed")
	default505Body = []byte("505 System Error")
	defaultConf    *EngConfig
)

type EngConfig struct {
	Network      string `yaml:"network"`
	Addr         string `yaml:"addr"`
	IdleTimeout  int64  `yaml:"idleTimeout"`
	ReadTimeout  int64  `yaml:"readTimeout"`
	WriteTimeout int64  `yaml:"writeTimeout"`
}

func init() {
	defaultConf = &EngConfig{}
	if err := config.Decode(config.CONFIG_FILE_NAME, "server", defaultConf); err != nil {
		panic("Service Load Config Fail " + err.Error())
	}
}

type HandlerFunc func(*Context)

type Engine struct {
	RouterGroup
	cfg       *EngConfig
	serve     atomic.Value
	wg        *sync.WaitGroup
	router    methodTrees
	allRouter map[string]map[string]struct{} // method -> path
	closed    int32                          // 1 为关闭 2为开启
	//Handlers []HandlerFunc
	noRouter []HandlerFunc
	noMethod []HandlerFunc

	reqChan chan *Context
}

//
//func buildDefaultConf() *EngConfig {
//	return &EngConfig{
//		Network:      "tcp",
//		Addr:         ":8086",
//		IdleTimeout:  1 * time.Hour,
//		ReadTimeout:  3 * time.Second,
//		WriteTimeout: 3 * time.Second,
//	}
//}

func NewEngine(cfg *EngConfig) *Engine {

	if cfg == nil {
		cfg = defaultConf
	}

	e := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			router:   true,
			BasePath: "/",
		},
		cfg:       cfg,
		wg:        &sync.WaitGroup{},
		router:    make(methodTrees, 0, 9),
		serve:     atomic.Value{},
		allRouter: make(map[string]map[string]struct{}),
		noRouter:  make([]HandlerFunc, 0),
		noMethod:  make([]HandlerFunc, 0),
		reqChan:   make(chan *Context, 1000), // 这里的1000 现在是随便给的
		closed:    CLOSED,
	}
	e.e = e
	// 加入pprof路由
	startPProf(e)
	e.addRouter("GET", "/allRouter", e.allRouters)
	return e
}

func (e *Engine) allRouters(c *Context) {
	if err := c.Json(e.allRouter, nil); err != nil {
		log.Error("all Routers Fail err is %s", err.Error())
	}
}

func (e *Engine) addRouter(method string, path string, handler ...HandlerFunc) {
	if path[0] != '/' {
		panic("path mast begin with '/'")
	}

	if method == "" {
		panic("http method can not empty ")
	}

	if len(handler) == 0 {
		panic("must be at least one handler")
	}

	allPath, ok := e.allRouter[method]
	if !ok {
		allPath = make(map[string]struct{})
		e.allRouter[method] = allPath
	}

	_, ok = allPath[path]

	if ok {
		panic(fmt.Sprintf("Path is Register method is %s , path is %s", method, path))
	}

	allPath[path] = struct{}{}

	router := e.router.getRoot(method)
	if router == nil {
		router = &node{}
		e.router = append(e.router, methodTree{root: router, method: method})
	}

	prelude := func(c *Context) {
		c.Method = method
		c.RouterPath = path
		log.Debug("router method is %s , path is %s ", method, path)
	}

	handlers := append([]HandlerFunc{prelude}, handler...)
	router.addRouter(path, handlers...)
	log.Info("Add Router  method is %s , path is %s", method, path)

	e.NoMethod(func(c *Context) {
		c.Byte(404, "text/plain; chatset=utf-8", default404Body)
		c.Abort()
	})

	e.NoRouter(func(c *Context) {
		c.Byte(405, "text/plain; chatset=utf-8", default405Body)
		c.Abort()
	})
}

func (e *Engine) NoRouter(handler ...HandlerFunc) {
	e.noRouter = handler
	e.rebuild404Handler()
}

func (e *Engine) NoMethod(handler ...HandlerFunc) {
	e.noMethod = handler
	e.rebuild405Handler()
}

func (e *Engine) rebuild404Handler() {
	e.noRouter = e.combineHandlers(e.noRouter...)
}

func (e *Engine) rebuild405Handler() {
	e.noMethod = e.combineHandlers(e.noMethod...)
}

func (e *Engine) isClosed() bool {

	if atomic.LoadInt32(&e.closed) == CLOSED {
		return true
	}
	return false
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if e.isClosed() {
		return
	}

	c := &Context{
		Res:     w,
		Req:     r,
		index:   -1,
		e:       e,
		handler: nil,
		Ctx:     nil,
	}

	e.HandlerContext(c)
	//
	//// 进行一个排队
	// TODO: 这里不能这样排队 应为过channel 会发生复制 放到channel后本次ServerHTTP就退出了 GC会把w也释放掉，所以会造成Write报错
	//select {
	//case e.reqChan <- c:
	//default:
	//}

}

//
//func (e *Engine) handler() {
//
//	for {
//		select {
//		case ctx, ok := <-e.reqChan:
//			if !ok {
//				e.wg.Done()
//				return
//			}
//			e.HandlerContext(ctx)
//		}
//
//	}
//
//}

func (e *Engine) HandlerContext(ctx *Context) {

	// parseFrom
	cType := ctx.Req.Header.Get("Content-Type")
	switch {
	case strings.Contains(cType, "multipart/form-data"):
		ctx.Req.ParseMultipartForm(MAX_MEM)
	default:
		ctx.Req.ParseForm()
	}
	var cancel context.CancelFunc
	ctx.Ctx, cancel = context.WithCancel(context.TODO())
	defer cancel()

	e.prepareHandler(ctx)
	ctx.Next()
}

func (e *Engine) prepareHandler(c *Context) {
	method := c.Req.Method
	path := c.Req.URL.Path

	h := e.router.getRoot(method).getHandler(path)
	if h == nil {
		c.handler = e.noRouter
		return
	}

	c.handler = h
	return
}

func (e *Engine) RunServer() {
	ser := &http.Server{
		Addr:         e.cfg.Addr,
		Handler:      e,
		IdleTimeout:  time.Duration(e.cfg.IdleTimeout) * time.Second,
		ReadTimeout:  time.Duration(e.cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(e.cfg.WriteTimeout) * time.Second,
	}

	e.serve.Store(ser)
	atomic.StoreInt32(&e.closed, START)

	//e.wg.Add(1)
	//go e.handler()

	e.wg.Add(1)
	go func() {
		if err := ser.ListenAndServe(); err != nil {
			log.Debug("RunServer ListenAndServer Err is %s ", err.Error())
		}
		e.wg.Done()
	}()
	log.Info("Http Server Start Addr is %s", e.cfg.Addr)
}

func (e *Engine) Close() {

	if e.isClosed() {
		return
	}

	atomic.StoreInt32(&e.closed, CLOSED)

	s := e.serve.Load().(*http.Server)
	s.Close()

	close(e.reqChan)

	e.wg.Wait()
	log.Info("Http Serve Stop Addr is %s", e.cfg.Addr)
}
