package binding

import "net/http"

type Binding interface {
	Name() string
	Bind(r *http.Request, data interface{}) error
	//testInterface(form map[string][]string, data interface{}) error
}

const (
	MEINJSON          = "application/json"
	MEINTEXT          = "text/plain"
	MEINPOSTFORM      = "application/x-www-urlencoded"
	METIMULTIPARTFORM = "multipart/form-data"
)

var (
	jsonBind          = jsonBinding{}
	queryBind         = queryBinding{}
	formBind          = formBinding{}
	postFormBind      = formPostBinding{}
	multipartFormBind = formMultipartBinding{}
)

func DefaultBind(method string, contentType string) Binding {
	if method == "GET" {
		//return formBind
		return queryBind
	}

	switch contentType {
	case MEINJSON:
		return jsonBind
	default:
		return formBind
	}
}
