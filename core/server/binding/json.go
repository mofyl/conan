package binding

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"net/http"
	"conan/utils"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type jsonBinding struct {
}

func (j jsonBinding) Name() string {
	return "json"
}

func (j jsonBinding) Bind(r *http.Request, data interface{}) error {
	decoder := utils.Json.NewDecoder(r.Body)

	if err := decoder.Decode(data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (j jsonBinding) testInterface(form map[string][]string, data interface{}) error {
	return nil
}
