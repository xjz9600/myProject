package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req        *http.Request
	Response   http.ResponseWriter
	PathParams map[string]string

	// 设置缓存
	cacheQueryValues url.Values
	MathRoute        string

	// 返回值
	RespStatusCode int
	RespData       []byte
	Tpl            TemplateEngine
	CacheSession   map[string]any
}

func (c *Context) BindJson(val any) error {
	return json.NewDecoder(c.Req.Body).Decode(val)
}

func (c *Context) FormValue(key string) (sv stringValue) {
	err := c.Req.ParseForm()
	if err != nil {
		sv.err = err
		return
	}
	values, ok := c.Req.PostForm[key]
	if !ok {
		sv.err = errors.New("web：没有对应的数据")
	}
	sv.val = values[0]
	return
}

func (c *Context) QueryValue(key string) (sv stringValue) {
	if c.cacheQueryValues == nil {
		c.cacheQueryValues = c.Req.URL.Query()
	}
	values, ok := c.cacheQueryValues[key]
	if !ok {
		sv.err = errors.New("web：连接中没有对应的请求参数")
	}
	sv.val = values[0]
	return
}

func (c *Context) ParamValue(key string) (sv stringValue) {
	val, ok := c.PathParams[key]
	if !ok {
		sv.err = errors.New("web：没有对应的路径参数")
	}
	sv.val = val
	return
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Response, cookie)
}

func (c *Context) RespJson(val any) {
	data, err := json.Marshal(val)
	if err != nil {
		c.RespStatusCode = http.StatusServiceUnavailable
		c.RespData = []byte(err.Error())
	}
	c.RespStatusCode = http.StatusOK
	c.RespData = data
}

func (c *Context) Render(templateName string, data any) {
	temp, err := c.Tpl.Render(c.Req.Context(), templateName, data)
	if err != nil {
		c.RespStatusCode = http.StatusInternalServerError
		c.RespData = []byte(err.Error())
		return
	}
	c.RespStatusCode = http.StatusOK
	c.RespData = temp
}

type stringValue struct {
	val string
	err error
}

func (s stringValue) AsInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}

func (s stringValue) AsString() (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return s.val, nil
}
