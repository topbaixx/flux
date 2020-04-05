package http

import (
	"context"
	"fmt"
	"github.com/bytepowered/flux"
	"github.com/bytepowered/flux/logger"
	"github.com/bytepowered/flux/pkg"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NewHttpExchange() *exchange {
	return &exchange{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

type exchange struct {
	httpClient *http.Client
}

func (e *exchange) Exchange(ctx flux.Context) *flux.InvokeError {
	ep := ctx.Endpoint()
	ret, err := e.Invoke(&ep, ctx)
	if err != nil {
		return err
	}
	resp := ret.(*http.Response)
	ctx.ResponseWriter().
		SetStatusCode(resp.StatusCode).
		SetHeaders(resp.Header).
		SetBody(resp.Body)
	return nil
}

func (e *exchange) Invoke(target *flux.Endpoint, ctx flux.Context) (interface{}, *flux.InvokeError) {
	httpRequest := ctx.RequestReader().HttpRequest()
	newRequest, err := e.Assemble(target, httpRequest)
	if nil != err {
		return nil, &flux.InvokeError{
			StatusCode: flux.StatusServerError,
			Message:    "HTTP:INVALID_REQUEST",
			Internal:   err,
		}
	} else {
		// Header透传以及传递AttrValues
		newRequest.Header = httpRequest.Header.Clone()
		for k, v := range ctx.AttrValues() {
			newRequest.Header.Set(k, pkg.ToString(v))
		}
	}
	resp, err := e.httpClient.Do(newRequest)
	if nil != err {
		msg := "HTTP:REMOTE_ERROR"
		if uErr, ok := err.(*url.Error); ok {
			msg = fmt.Sprintf("HTTP:REMOTE_ERROR:%s", uErr.Error())
		}
		return nil, &flux.InvokeError{
			StatusCode: flux.StatusServerError,
			Message:    msg,
			Internal:   err,
		}
	}
	return resp, nil
}

func (e *exchange) Assemble(endpoint *flux.Endpoint, inReq *http.Request) (*http.Request, error) {
	inParams := endpoint.Arguments
	newQuery := inReq.URL.RawQuery
	// 使用可重复读的GetBody函数
	reader, err := inReq.GetBody()
	if nil != err {
		return nil, fmt.Errorf("get body by func, err: %w", err)
	}
	defer pkg.CloseSilently(reader)
	var newBodyReader io.Reader = reader
	if len(inParams) > 0 {
		// 如果Endpoint定义了参数，即表示限定参数传递
		data := _toHttpUrlValues(inParams).Encode()
		// GET：参数拼接到URL中；
		if http.MethodGet == endpoint.UpstreamMethod {
			if newQuery == "" {
				newQuery = data
			} else {
				newQuery += "&" + data
			}
		} else {
			// 其它方法：拼接到Body中，并设置form-data/x-www-url-encoded
			newBodyReader = strings.NewReader(data)
		}
	}
	// 未定义参数，即透传Http请求：Rewrite inReq path
	newUrl := &url.URL{
		Host:       endpoint.UpstreamHost,
		Path:       endpoint.UpstreamUri,
		Scheme:     inReq.URL.Scheme,
		Opaque:     inReq.URL.Opaque,
		User:       inReq.URL.User,
		RawPath:    inReq.URL.RawPath,
		ForceQuery: inReq.URL.ForceQuery,
		RawQuery:   newQuery,
		Fragment:   inReq.URL.Fragment,
	}
	timeout, err := time.ParseDuration(endpoint.RpcTimeout)
	if err != nil {
		logger.Warnf("Illegal endpoint rpc-timeout: ", endpoint.RpcTimeout)
		timeout = time.Second * 10
	}
	stdCtx, _ := context.WithTimeout(context.Background(), timeout)
	newRequest, err := http.NewRequestWithContext(stdCtx, endpoint.UpstreamMethod, newUrl.String(), newBodyReader)
	if nil != err {
		return nil, fmt.Errorf("new request, method: %s, url: %s, err: %w", endpoint.UpstreamMethod, newUrl, err)
	}
	// Body数据设置application/x-www-url-encoded
	if http.MethodGet != endpoint.UpstreamMethod {
		newRequest.Header.Set("Content-TypeId", "application/x-www-form-urlencoded")
	}
	newRequest.Header.Set("User-Agent", "FluxGo/Exchange/v1")
	return newRequest, err
}
