package server

import (
	"conan/core/server/binding"
	"conan/core/server/rending"
	"context"
	"fmt"
	"math"
	"net/http"
	"runtime"
	//"conan/ecode"
	"conan/log"
)

const (
	ABORT_INDEX = math.MaxInt8 / 2
)

type Chain func(c *Context)

type Context struct {
	Req *http.Request
	Res http.ResponseWriter
	Ctx context.Context

	e       *Engine
	index   int8
	handler []HandlerFunc

	Method     string
	RouterPath string
	//Err        ecode.ErrMsgs
}

func (c *Context) Abort() {
	c.index = ABORT_INDEX
}

func (c *Context) Status(code int) {
	c.Res.WriteHeader(code)
}

func (c *Context) Next() {

	defer func() {

		if err := recover(); err != nil {
			runtime.Caller(1)
			buf := make([]byte, 4086)
			n := runtime.Stack(buf, false)
			pl := fmt.Sprintf("http server panic: %v\n%s\n", err, buf[:n])
			log.Error("http server recover  is %s", pl)
			c.Byte(500, "text/plain", default505Body)
		}

	}()

	c.index++
	// 计算时间
	for c.handler != nil && c.index < int8(len(c.handler)) {
		c.handler[c.index](c)
		c.index++
	}
	// 采集rtt
}

func (c *Context) Render(code int, render rending.Render) error {
	render.WriteContentType(c.Res)

	if code != http.StatusOK {
		c.Status(code)
	}
	//if c.Err != nil {
	//
	//}
	err := render.Render(c.Res)
	if err != nil {
		log.Error("Context: Write Response Err is %s ", err.Error())
		return err
	}
	c.Abort()
	if c.Req.Body != nil {
		c.Req.Body.Close()
	}
	return nil
}

func (c *Context) Json(data interface{}, err error) error {
	code := http.StatusOK

	//var eCode ecode.ErrMsgs
	//var ok bool
	//if err == nil {
	//	eCode = ecode.OK
	//	ok = true
	//} else {
	//	eCode, ok = err.(ecode.ErrMsgs)
	//	if !ok {
	//		log.Error("Context Json Fail Convert Err Fail")
	//		return errors.New("Context Json Fail Convert Err Fail")
	//	}
	//
	//}

	//c.Err = eCode
	return c.Render(code, &rending.Json{
		//Code: eCode.Code(),
		//Msg:  eCode.Error(),
		Data: data,
	})

}

func (c *Context) String(code int, format string, data ...interface{}) error {
	return c.Render(code, &rending.String{
		Format: format,
		Data:   data,
	})
}

func (c *Context) Byte(code int, contentType string, data ...[]byte) error {
	return c.Render(code, &rending.Data{
		ContentType: contentType,
		Data:        data,
	})
}

func (c *Context) Bind(obj interface{}) error {
	bind := binding.DefaultBind(c.Req.Method, c.Req.Header.Get("Content-Type"))
	return c.mustBind(bind, obj)
}

func (c *Context) mustBind(bind binding.Binding, obj interface{}) error {
	return bind.Bind(c.Req, obj)
}
