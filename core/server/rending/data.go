package rending

import (
	"github.com/pkg/errors"
	"net/http"
)

type Data struct {
	ContentType string
	Data [][]byte
}

func (d *Data)WriteContentType(w http.ResponseWriter) {
	writeContentType(w , []string{d.ContentType})
}
func (d *Data)Render(w http.ResponseWriter) error {
	writeContentType(w , []string{d.ContentType})
	var err error
	for _ , v := range d.Data{
		if _ , err := w.Write(v); err != nil {
			err = errors.WithStack(err)
		}
	}
	return err
}


