package server

import (
	"github.com/bytepowered/flux/ext"
	"github.com/bytepowered/flux/internal"
	"net/http"
	"strings"
)

const (
	_typeApplication = "application"
	_typeProtocol    = "protocol"
	_typeHttpPattern = "http-pattern"
	_typeUpstreamUri = "upstream-uri"
)

type _filter func(ep *internal.MultiVersionEndpoint) bool

// 支持以下过滤条件
var _typeKeys = []string{"application", "protocol", "http-pattern", "upstream-uri"}

var (
	_filterFactories = make(map[string]func(string) _filter)
)

func DebugQueryEndpoint(datamap map[string]*internal.MultiVersionEndpoint) http.HandlerFunc {
	// Endpoint查询
	json := ext.GetSerializer(ext.TypeNameSerializerJson)
	return func(writer http.ResponseWriter, request *http.Request) {
		if data, err := json.Marshal(queryEndpoints(datamap, request)); nil != err {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte(err.Error()))
		} else {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write(data)
		}
	}
}

func init() {
	_filterFactories[_typeApplication] = func(query string) _filter {
		return func(ep *internal.MultiVersionEndpoint) bool {
			return query == ep.RandomVersion().Application
		}
	}
	_filterFactories[_typeProtocol] = func(query string) _filter {
		return func(ep *internal.MultiVersionEndpoint) bool {
			return strings.ToLower(query) == strings.ToLower(ep.RandomVersion().Protocol)
		}
	}
	_filterFactories[_typeHttpPattern] = func(query string) _filter {
		return func(ep *internal.MultiVersionEndpoint) bool {
			return query == ep.RandomVersion().HttpPattern
		}
	}
	_filterFactories[_typeUpstreamUri] = func(query string) _filter {
		return func(ep *internal.MultiVersionEndpoint) bool {
			return query == ep.RandomVersion().UpstreamUri
		}
	}
}

func queryEndpoints(data map[string]*internal.MultiVersionEndpoint, request *http.Request) interface{} {
	filters := make([]_filter, 0)
	query := request.URL.Query()
	for _, key := range _typeKeys {
		if query := query.Get(key); "" != query {
			if f, ok := _filterFactories[key]; ok {
				filters = append(filters, f(query))
			}
		}
	}
	if len(filters) == 0 {
		m := make(map[string]interface{})
		for k, v := range data {
			m[k] = v.ToSerializableMap()
		}
		return m
	}
	return _queryWithFilters(data, filters...)
}

func _queryWithFilters(data map[string]*internal.MultiVersionEndpoint, filters ..._filter) []interface{} {
	items := make([]interface{}, 0)
DataLoop:
	for _, v := range data {
		for _, filter := range filters {
			// 任意Filter返回True
			if filter(v) {
				items = append(items, v.ToSerializableMap())
				continue DataLoop
			}
		}
	}
	return items
}