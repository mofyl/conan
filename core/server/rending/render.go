package rending

import (
	"net/http"
)

type RenderType string

var (
	JSON   RenderType = "JSON"
	BYTE   RenderType = "BYTE"
	STRING RenderType = "STRING"
)

type Base struct {
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

type Render interface {
	WriteContentType(w http.ResponseWriter)
	Render(w http.ResponseWriter) error
}

func writeContentType(w http.ResponseWriter, value []string) {
	head := w.Header()
	if val := head["Content-Type"]; len(val) == 0 {
		head["Content-Type"] = value
	}
}
