package dubbo

import (
	"github.com/apache/dubbo-go/protocol"
	"github.com/bytepowered/flux/flux-node"
	"github.com/spf13/cast"
	"net/http"
)

const (
	ResponseKeyStatusCode = "@net.bytepowered.flux.http-status"
	ResponseKeyHeaders    = "@net.bytepowered.flux.http-headers"
)

func NewTransportCodecFuncWith(codeKey, headerKey string) flux.TransportCodec {
	return func(ctx *flux.Context, raw interface{}) (*flux.ResponseBody, error) {
		// 支持Dubbo返回Result类型
		rpcr, ok := raw.(protocol.Result)
		if !ok {
			return &flux.ResponseBody{
				StatusCode: flux.StatusOK, Headers: make(http.Header, 0), Body: raw,
			}, nil
		}
		attrs := make(map[string]interface{}, 8)
		if err := rpcr.Error(); nil != err {
			return nil, err
		}
		data := rpcr.Result()
		status := flux.StatusOK
		for k, v := range rpcr.Attachments() {
			if k == codeKey {
				status = cast.ToInt(v)
			} else if k == headerKey {
				// TODO 需要更新Attachment类型为map[string]interface{}
			} else {
				attrs[k] = v
			}
		}
		return &flux.ResponseBody{
			StatusCode: status, Headers: make(http.Header, 0), Attachments: attrs, Body: data,
		}, nil
	}
}

func NewTransportCodecFunc() flux.TransportCodec {
	return NewTransportCodecFuncWith(ResponseKeyStatusCode, ResponseKeyHeaders)
}
