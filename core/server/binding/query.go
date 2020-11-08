package binding

import (
	"github.com/pkg/errors"
	"net/http"
)

type queryBinding struct{}

func (queryBinding) Name() string {
	return "query"
}

func (queryBinding) Bind(r *http.Request, data interface{}) error {
	if err := mapForm(data, r.URL.Query()); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (queryBinding) testInterface(form map[string][]string, data interface{}) error {
	return nil
}
