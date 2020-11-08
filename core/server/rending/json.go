package rending

import (
	"net/http"
	"conan/log"
	"conan/utils"
)

var (
	jsonContentType = []string{"application/json; chatset=utf-8"}
)

type JsonEncode interface {
	ToJson() []byte
}

type Json struct {
	Code int64       `json:"code,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func (j *Json) ToJson() ([]byte, error) {
	return utils.Json.Marshal(j)
}

func (j *Json) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func (j *Json) Render(w http.ResponseWriter) error {
	writeContentType(w, jsonContentType)

	b, err := utils.Json.Marshal(j)

	if err != nil {
		log.Error("WriteResponse msgPack Marshal Err is %s ", err.Error())
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}
