package server

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"
	"conan/core/server/rending"
	//"conan/ecode"
)

var (
	cli *http.Client
)

func init() {
	cli = &http.Client{}
	cli.Timeout = 5 * time.Second
}

//func HttpCode(c *Context, msgs ecode.ErrMsgs) {
//	c.Json(nil, msgs)
//}
//
//func HttpData(c *Context, data interface{}) {
//	c.Json(data, ecode.OK)
//}
//
//func BackCode(urlStr string, code ecode.ErrMsgs) (*http.Response, error) {
//	return HttpJsonPost(urlStr, &rending.Json{
//		Code: code.Code(),
//		Msg:  code.Error(),
//		Data: nil,
//	})
//}
//
//func BackData(urlStr string, data interface{}) (*http.Response, error) {
//	return HttpJsonPost(urlStr, &rending.Json{
//		Code: ecode.OK.Code(),
//		Msg:  ecode.OK.Error(),
//		Data: data,
//	})
//}

func HttpJsonPost(urlStr string, data *rending.Json) (*http.Response, error) {

	b, _ := data.ToJson()

	return cli.Post(urlStr, "application/json", bytes.NewReader(b))
}

func HttpGet(url string, arg map[string]string) (*http.Response, error) {

	builder := strings.Builder{}

	builder.WriteString(url)

	builder.WriteString("?")
	for k, v := range arg {
		builder.WriteString(fmt.Sprintf("%s=%s&", k, v))
	}

	urlStr := builder.String()
	urlStr = urlStr[:len(urlStr)-1]

	return cli.Get(urlStr)
}

func Close() {
	if cli != nil {
		cli.CloseIdleConnections()
	}
}
