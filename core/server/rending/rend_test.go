package rending

import (
	"encoding/json"
	"net/http"
	"testing"
	"conan/utils"
)

func (b *Base) ToJson() ([]byte, error) {
	return utils.Json.Marshal(b)
}

type Rend struct {
	Base
	A string `json:"a"`
	B string `json:"b"`
}

func (r *Rend) ServeHTTP(w http.ResponseWriter, request *http.Request) {

	r.A = "zzz"
	r.B = "xccxvcxv"
	r.ErrCode = 1231
	r.ErrMsg = "cvbvcb"

	b, _ := json.Marshal(r)
	w.WriteHeader(200)
	w.Write(b)
}

func (r *Rend) ToJson() ([]byte, error) {
	return utils.Json.Marshal(r)
}

func TestJson(t *testing.T) {
	r := &Rend{}
	http.ListenAndServe(":9000", r)
}
