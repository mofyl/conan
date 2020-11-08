package rending

import (
	"fmt"
	"io"
	"net/http"
)

var (
	strContentType = []string{"text/plain; chatset=utf-8"}
)

type String struct {
	Data []interface{}
	Format string
}
func (s *String)WriteContentType(w http.ResponseWriter) {
	writeContentType(w, strContentType)
}
func (s *String)Render(w http.ResponseWriter) error {
	writeContentType(w , strContentType)
	var err error
	if len(s.Data) > 0 {
		_ , err = fmt.Fprintf(w , s.Format , s.Data)
	}else {
		_ , err = io.WriteString(w , s.Format)
	}
	return err
}




