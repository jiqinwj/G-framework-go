package pkg

import (
	"fmt"
	"net/http"
)

// 这个写法可用于确保 HTTPHandlerMap 实现 HTTPHandler 接口。
// 如果 HTTPHandlerMap 没有实现 HTTPHandler 接口，这里就会报错。
var _ HTTPHandler = &HTTPHandlerMap{}

// HTTPHandlerMap 基于 map 实现路由处理
type HTTPHandlerMap struct {
	// mapRoute 路由对应的处理方法
	mapRoute map[string]HTTPHandlerFunc
}

func NewHTTPHandlerMap() HTTPHandler {
	return &HTTPHandlerMap{
		mapRoute: make(map[string]HTTPHandlerFunc),
	}
}

// HandlerHTTP HTTPHandler.HandlerHTTP
func (p1h *HTTPHandlerMap) HandlerHTTP(c *HTTPContext) {
	p1req := c.P1req
	fmt.Printf("HTTPHandlerMap, HandlerHTTP, p1req.Method: %s, p1req.URL.Path: %s\n", p1req.Method, p1req.URL.Path)

	// 路由查询，找到对应的处理函数
	hhFunc, ok := p1h.mapRoute[p1req.Method+"#"+p1req.URL.Path]
	if !ok {
		c.P1resW.WriteHeader(http.StatusNotFound)
		_, _ = c.P1resW.Write([]byte("route not found"))
		return
	}

	hhFunc(c)
}

// RegisteRoute HTTPHandler.HTTPRoute.RegisteRoute
func (p1h *HTTPHandlerMap) RegisteRoute(method string, pattern string, hhFunc HTTPHandlerFunc) error {
	// 这里用 HTTP 方法和路由构造一个唯一键，实现区分不同 HTTP 方法的路由
	p1h.mapRoute[method+"#"+pattern] = hhFunc
	return nil
}
