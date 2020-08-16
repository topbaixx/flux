package echo

import (
	"context"
	"github.com/bytepowered/flux/webx"
	"github.com/labstack/echo/v4"
)

var _ webx.WebServer = new(AdaptWebServer)

func NewAdaptWebServer() webx.WebServer {
	server := echo.New()
	server.HideBanner = true
	server.HidePort = true
	// 可重读Body
	server.Pre(RepeatableReadBody)
	return &AdaptWebServer{server}
}

type AdaptWebServer struct {
	server *echo.Echo
}

func (w *AdaptWebServer) SetRouteNotFoundHandler(fun webx.WebRouteHandler) {
	echo.NotFoundHandler = AdaptWebRouteHandler(fun).AdaptFunc
}

func (w *AdaptWebServer) SetWebErrorHandler(fun webx.WebErrorHandler) {
	w.server.HTTPErrorHandler = func(err error, c echo.Context) {
		fun(err, toAdaptWebContext(c))
	}
}

func (w *AdaptWebServer) AddWebInterceptor(m webx.WebMiddleware) {
	w.server.Pre(AdaptWebMiddleware(m).AdaptFunc)
}

func (w *AdaptWebServer) AddWebMiddleware(m webx.WebMiddleware) {
	w.server.Use(AdaptWebMiddleware(m).AdaptFunc)
}

func (w *AdaptWebServer) AddWebRouteHandler(method, pattern string, h webx.WebRouteHandler, m ...webx.WebMiddleware) {
	wms := make([]echo.MiddlewareFunc, len(m))
	for i, mi := range m {
		wms[i] = AdaptWebMiddleware(mi).AdaptFunc
	}
	w.server.Add(method, pattern, AdaptWebRouteHandler(h).AdaptFunc, wms...)
}

func (w *AdaptWebServer) Start(addr string) error {
	return w.server.Start(addr)
}

func (w *AdaptWebServer) StartTLS(addr string, certFile, keyFile string) error {
	return w.server.StartTLS(addr, certFile, keyFile)
}

func (w *AdaptWebServer) Shutdown(ctx context.Context) error {
	return w.server.Shutdown(ctx)
}