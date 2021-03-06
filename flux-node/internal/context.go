package internal

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bytepowered/flux/flux-node"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/url"
)

var _ flux.ServerWebContext = new(EchoWebContext)

func NewServeWebContext(ctx echo.Context, reqid string, listener flux.WebListener) flux.ServerWebContext {
	return &EchoWebContext{
		echoc:     ctx,
		listener:  listener,
		context:   context.WithValue(ctx.Request().Context(), keyRequestId, reqid),
		variables: make(map[interface{}]interface{}, 16),
	}
}

type EchoWebContext struct {
	listener  flux.WebListener
	context   context.Context
	echoc     echo.Context
	variables map[interface{}]interface{}
}

func (w *EchoWebContext) WebListener() flux.WebListener {
	return w.listener
}

func (w *EchoWebContext) ShadowContext() echo.Context {
	return w.echoc
}

func (w *EchoWebContext) RequestId() string {
	return w.context.Value(keyRequestId).(string)
}

func (w *EchoWebContext) Context() context.Context {
	return w.context
}

func (w *EchoWebContext) Request() *http.Request {
	return w.echoc.Request()
}

func (w *EchoWebContext) URI() string {
	return w.Request().RequestURI
}

func (w *EchoWebContext) URL() *url.URL {
	return w.Request().URL
}

func (w *EchoWebContext) Method() string {
	return w.Request().Method
}

func (w *EchoWebContext) Host() string {
	return w.Request().Host
}

func (w *EchoWebContext) RemoteAddr() string {
	return w.Request().RemoteAddr
}

func (w *EchoWebContext) HeaderVars() http.Header {
	return w.Request().Header
}

func (w *EchoWebContext) QueryVars() url.Values {
	return w.echoc.QueryParams()
}

func (w *EchoWebContext) PathVars() url.Values {
	names := w.echoc.ParamNames()
	copied := make(url.Values, len(names))
	for _, n := range names {
		copied.Set(n, w.echoc.Param(n))
	}
	return copied
}

func (w *EchoWebContext) FormVars() url.Values {
	f, _ := w.echoc.FormParams()
	return f
}

func (w *EchoWebContext) CookieVars() []*http.Cookie {
	return w.echoc.Cookies()
}

func (w *EchoWebContext) HeaderVar(name string) string {
	return w.Request().Header.Get(name)
}

func (w *EchoWebContext) QueryVar(name string) string {
	return w.echoc.QueryParam(name)
}

func (w *EchoWebContext) PathVar(name string) string {
	// use cached vars
	return w.echoc.Param(name)
}

func (w *EchoWebContext) FormVar(name string) string {
	return w.echoc.FormValue(name)
}

func (w *EchoWebContext) CookieVar(name string) (*http.Cookie, error) {
	return w.echoc.Cookie(name)
}

func (w *EchoWebContext) BodyReader() (io.ReadCloser, error) {
	return w.Request().GetBody()
}

func (w *EchoWebContext) Rewrite(method string, path string) {
	if "" != method {
		w.Request().Method = method
	}
	if "" != path {
		w.Request().URL.Path = path
	}
}

func (w *EchoWebContext) Write(statusCode int, contentType string, data []byte) error {
	return w.WriteStream(statusCode, contentType, bytes.NewReader(data))
}

func (w *EchoWebContext) WriteStream(statusCode int, contentType string, reader io.Reader) error {
	writer := w.echoc.Response()
	writer.Header().Set(echo.HeaderContentType, contentType)
	writer.WriteHeader(statusCode)
	if _, err := io.Copy(writer, reader); nil != err {
		return fmt.Errorf("web context write failed, error: %w", err)
	}
	return nil
}

func (w *EchoWebContext) SetResponseWriter(rw http.ResponseWriter) {
	w.echoc.Response().Writer = rw
}

func (w *EchoWebContext) ResponseWriter() http.ResponseWriter {
	return w.echoc.Response().Writer
}

func (w *EchoWebContext) Variable(key string) interface{} {
	v, _ := w.GetVariable(key)
	return v
}

func (w *EchoWebContext) SetVariable(key string, value interface{}) {
	w.variables[key] = value
}

func (w *EchoWebContext) GetVariable(key string) (interface{}, bool) {
	// 本地Variable
	v, ok := w.variables[key]
	if ok {
		return v, true
	}
	// 从Context中加载
	v = w.echoc.Get(key)
	return v, nil != v
}
