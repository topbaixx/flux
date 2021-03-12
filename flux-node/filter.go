package flux

import (
	"fmt"
	"github.com/spf13/cast"
	"net/http"
)

// ServeError 定义网关处理请求的服务错误；
// 它包含：错误定义的状态码、错误消息、内部错误等元数据
type ServeError struct {
	StatusCode int                    // 响应状态码
	ErrorCode  interface{}            // 业务错误码
	Message    string                 // 错误消息
	CauseError error                  // 内部错误对象；错误对象不会被输出到请求端；
	Header     http.Header            // 响应Header
	Extras     map[string]interface{} // 用于定义和跟踪的额外信息；额外信息不会被输出到请求端；
}

func (e *ServeError) Error() string {
	if nil != e.CauseError {
		return fmt.Sprintf("ServeError: StatusCode=%d, ErrorCode=%s, Message=%s, Extras=%+v, Error=%s", e.StatusCode, e.ErrorCode, e.Message, e.Extras, e.CauseError)
	} else {
		return fmt.Sprintf("ServeError: StatusCode=%d, ErrorCode=%s, Message=%s, Extras=%+v", e.StatusCode, e.ErrorCode, e.Message, e.Extras)
	}
}

func (e *ServeError) GetErrorCode() string {
	return cast.ToString(e.ErrorCode)
}

func (e *ServeError) GetExtras(key string) interface{} {
	return e.Extras[key]
}

func (e *ServeError) SetExtras(key string, value interface{}) {
	if e.Extras == nil {
		e.Extras = make(map[string]interface{}, 4)
	}
	e.Extras[key] = value
}

func (e *ServeError) MergeHeader(header http.Header) *ServeError {
	if e.Header == nil {
		e.Header = header.Clone()
	} else {
		for key, values := range header {
			for _, value := range values {
				e.Header.Add(key, value)
			}
		}
	}
	return e
}

type (
	// FilterHandler 定义一个处理方法，处理请求Context；如果发生错误则返回 ServeError。
	FilterHandler func(Context) *ServeError
	// FilterSkipper 定义一个函数，用于Filter执行中跳过某些处理。返回True跳过某些处理，见具体Filter的实现逻辑。
	FilterSkipper func(Context) bool
	// Filter 用于定义处理方法的顺序及内在业务逻辑
	Filter interface {
		// TypeId Filter的类型标识
		TypeId() string
		// DoFilter 执行Filter链
		DoFilter(next FilterHandler) FilterHandler
	}
)

type (
	// FilterSelector 用于请求处理前的动态选择Filter
	FilterSelector interface {
		// Activate 返回当前请求是否激活Selector
		Activate(ctx Context) bool
		// DoSelect 根据请求返回激活的选择器列表
		DoSelect(ctx Context) []Filter
	}
)