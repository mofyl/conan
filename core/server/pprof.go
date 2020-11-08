package server

import (
	"net/http"
	"net/http/pprof"
)

func startPProf(e *Engine) {

	/*


		Count	Profile
		1	allocs
		0	block
		0	cmdline
		22	goroutine
		1	heap
		0	mutex
		0	profile
		15	threadcreate
		0	trace


	*/

	e.GET("/debug/pprof/", pprofHandler(pprof.Index))

	e.GET("/debug/pprof/allocs", pprofHandler(pprof.Handler("allocs").ServeHTTP))
	e.GET("/debug/pprof/block", pprofHandler(pprof.Handler("block").ServeHTTP))
	e.GET("/debug/pprof/cmdline", pprofHandler(pprof.Cmdline))
	e.GET("/debug/pprof/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
	e.GET("/debug/pprof/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
	e.GET("/debug/pprof/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
	e.GET("/debug/pprof/profile", pprofHandler(pprof.Profile))
	e.GET("/debug/pprof/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
	e.GET("/debug/pprof/trace", pprofHandler(pprof.Trace))
}

func pprofHandler(h http.HandlerFunc) HandlerFunc {
	return func(c *Context) {
		h.ServeHTTP(c.Res, c.Req)
	}
}
