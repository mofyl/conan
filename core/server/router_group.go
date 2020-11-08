package server

import "conan/utils"

type IRouter interface {
	UseFunc(middle ...HandlerFunc) IRouter

	PUT(router string, handlerFunc ...HandlerFunc) IRouter
	POST(router string, handlerFunc ...HandlerFunc) IRouter
	GET(router string, handlerFunc ...HandlerFunc) IRouter
	DELETE(router string, handlerFunc ...HandlerFunc) IRouter
	HEAD(router string, handlerFunc ...HandlerFunc) IRouter
}

type RouterGroup struct {
	Handlers []HandlerFunc
	router   bool
	BasePath string
	e        *Engine
}

func (r *RouterGroup) NewGroup(basePath string, handlers ...HandlerFunc) IRouter {
	return &RouterGroup{
		Handlers: r.combineHandlers(handlers...),
		router:   false,
		BasePath: r.calcAbsPath(basePath),
		e:        r.e,
	}
}

func (r *RouterGroup) combineHandlers(handler ...HandlerFunc) []HandlerFunc {
	finalSize := len(r.Handlers) + len(handler)

	mergeHandler := make([]HandlerFunc, 0, finalSize)

	mergeHandler = append(mergeHandler, r.Handlers...)
	mergeHandler = append(mergeHandler, handler...)
	return mergeHandler
}

func (r *RouterGroup) UseFunc(middle ...HandlerFunc) IRouter {
	r.Handlers = append(r.Handlers, middle...)
	return r
}

func (r *RouterGroup) calcAbsPath(relativePath string) string {
	return utils.JoinPath(r.BasePath, relativePath)
}

func (r *RouterGroup) hand(method string, router string, handlers ...HandlerFunc) *RouterGroup {
	finalPath := r.calcAbsPath(router)

	finalHandlers := r.combineHandlers(handlers...)

	r.e.addRouter(method, finalPath, finalHandlers...)
	return r
}

func (r *RouterGroup) POST(router string, handlers ...HandlerFunc) IRouter {
	r.hand("POST", router, handlers...)
	return r
}

func (r *RouterGroup) GET(router string, handlers ...HandlerFunc) IRouter {
	r.hand("GET", router, handlers...)
	return r
}

func (r *RouterGroup) PUT(router string, handlers ...HandlerFunc) IRouter {
	r.hand("PUT", router, handlers...)
	return r
}

func (r *RouterGroup) DELETE(router string, handlers ...HandlerFunc) IRouter {
	r.hand("DELETE", router, handlers...)
	return r
}

func (r *RouterGroup) HEAD(router string, handlers ...HandlerFunc) IRouter {
	r.hand("HEAD", router, handlers...)
	return r
}
