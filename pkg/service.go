package pkg

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

// Service 服务接口
type Service interface {
	Start(addr string, port string) error
	// Shutdown 服务关闭
	Shutdown(c context.Context) error
	HTTPRoute
}

// HTTPService HTTP 服务
type HTTPService struct {
	Name     string
	handler  HTTPHandler
	entrance MiddlewareFunc
	// 资源池，复用 HTTPContext
	hcPool sync.Pool
	// 启动服务时，保存实例，用于服务关闭
	p1HTTPServer *http.Server
}

// NewHTTPSrevice 创建一个 Service 接口的实例，指定服务的名字和中间件组
func NewHTTPSrevice(name string, arr1Builder ...MiddlewareBuilder) Service {
	var p1h HTTPHandler = NewHTTPHandlerTree()
	var hf MiddlewareFunc = p1h.HandlerHTTP
	for i := len(arr1Builder) - 1; i > -1; i-- {
		var mf MiddlewareBuilder = arr1Builder[i]
		hf = mf(hf)
	}

	return &HTTPService{
		Name:     name,
		handler:  p1h,
		entrance: hf,
		hcPool: sync.Pool{New: func() interface{} {
			return NewHTTPContext()
		}},
	}
}

// ServeHTTP Handler.ServeHTTP
func (p1s *HTTPService) ServeHTTP(p1resW http.ResponseWriter, p1req *http.Request) {
	// 从资源池里获取 HTTPContext 实例
	p1c := p1s.hcPool.Get().(*HTTPContext)
	// 归还资源到资源池
	defer func() {
		p1s.hcPool.Put(p1c)
	}()
	p1c.Reset(p1resW, p1req)
	p1s.entrance(p1c)
}

// Start Service.Start
func (p1s *HTTPService) Start(addr string, port string) error {
	fmt.Printf("HTTPService %s start at %s...\r\n", p1s.Name, addr+":"+port)
	p1s.p1HTTPServer = &http.Server{
		Addr:    addr + ":" + port,
		Handler: p1s,
	}
	return p1s.p1HTTPServer.ListenAndServe()
}

// Shutdown Service.Shutdown
func (p1s *HTTPService) Shutdown(c context.Context) error {
	fmt.Printf("HTTPService %s shutdown\r\n", p1s.Name)
	return p1s.p1HTTPServer.Shutdown(c)
}

// RegisteRoute Service.HTTPRoute.RegisteRoute
func (p1s *HTTPService) RegisteRoute(method string, pattern string, hhFunc HTTPHandlerFunc) error {
	return p1s.handler.RegisteRoute(method, pattern, hhFunc)
}
