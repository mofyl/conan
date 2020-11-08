package binding

import (
	"net/http"
	"reflect"
	"sync"
)

var (
	sCache = &cache{
		data:  make(map[reflect.Type]*sInfo),
		mutex: sync.RWMutex{},
	}
)

type cache struct {
	data  map[reflect.Type]*sInfo
	mutex sync.RWMutex
}

type sInfo struct {
	fields []*field
}

type option map[string]struct{}

type field struct {
	tp      reflect.StructField
	name    string
	options option

	hasDefault   bool
	defaultValue reflect.Value
}

func (c *cache) get(obj reflect.Type) *sInfo {
	c.mutex.RLock()
	var s *sInfo
	var ok bool
	if s, ok = c.data[obj]; !ok {
		c.mutex.RUnlock()
		// 解析并缓存该 type
		s = c.set(obj)
		return s
	}
	c.mutex.RUnlock()
	return s
}

func (c *cache) set(p reflect.Type) *sInfo {
	tp := p.Elem()
	s := &sInfo{
		fields: make([]*field, 0, tp.NumField()),
	}
	for i := 0; i < tp.NumField(); i++ {
		fd := &field{}
		fd.tp = tp.Field(i)
		//fd.name = fd.tp.Tag.Get("form")
		info := fd.tp.Tag.Get("form")
		name, op := splitNameAndOption(info)
		fd.name = name
		fd.options = op
		if dev := fd.tp.Tag.Get("default"); dev != "" {
			dv := reflect.New(fd.tp.Type).Elem()
			err := setWithProperType(fd.tp.Type.Kind(), []string{dev}, dv, fd.options)
			if err != nil {
				continue
			}
			fd.hasDefault = true
			fd.defaultValue = dv
		}
		s.fields = append(s.fields, fd)
	}

	c.mutex.Lock()
	c.data[p] = s
	c.mutex.Unlock()
	return s
}

// ================== 分割线 ========================

const (
	defaultMem = 32 * 1024 * 1024
)

type formBinding struct{}

func (formBinding) Name() string {
	return "form"
}

func (formBinding) Bind(r *http.Request, data interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return mapForm(data, r.Form)
}

func (formBinding) testInterface(form map[string][]string, data interface{}) error {
	return mapForm(data, form)
}

type formPostBinding struct{}

func (formPostBinding) Name() string {
	return "form-urlencoded" // 这类型的form会将 form中的内容转换为键值对
}

func (formPostBinding) Bind(r *http.Request, data interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	return mapForm(data, r.PostForm)
}

func (formPostBinding) testInterface(form map[string][]string, data interface{}) error {
	return mapForm(data, form)
}

type formMultipartBinding struct{}

func (formMultipartBinding) Name() string {
	return "multipart/form-data" // 一般使用这类的上传文件
}

func (formMultipartBinding) Bind(r *http.Request, data interface{}) error {
	if err := r.ParseMultipartForm(defaultMem); err != nil {
		return err
	}

	return mapForm(data, r.MultipartForm.Value)
}

func (formMultipartBinding) testInterface(form map[string][]string, data interface{}) error {
	return mapForm(data, form)
}
