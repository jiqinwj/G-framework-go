package pkg

import (
	"net/http"
	"sync"
)

// 路由对应的处理方法的定义
type HTTPHandlerFunc func(p7ctx *HTTPContext)

// 核心处理逻辑的接口定义
type HTTPHandlerInterface interface {
	http.Handler
	RouteInterface
	MiddlewareInterface
}

// HTTPHandler 核心处理逻辑
type HTTPHandler struct {
	router
	// 全局中间件
	s5f4middleware []HTTPMiddleware
	// 内存池 复用 HTTPContext
	ctxPool sync.Pool
	// isRunning 服务是否正在运行
	isRunning bool
}

func (p7this *HTTPHandler) ServeHTTP(i9w http.ResponseWriter, p7r *http.Request) {
	// 从内存池里获取 HTTPContext
	p7ctx := p7this.ctxPool.Get().(*HTTPContext)
	// 归还资源到资源池
	defer func() {
		p7ctx.R
	}()

}

//var _ HTTPHandlerInterface = &HTTPHandler{}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{
		router: newRouter(),
		ctxPool: sync.Pool{
			New: func() interface{} {
				return NewHTTPContext()
			},
		},
		isRunning: true,
	}
}

// 增加中间件
func (p7this *HTTPHandler) AddMiddleware(s5f4mw ...HTTPMiddleware) {
	if nil == p7this.s5f4middleware {
		p7this.s5f4middleware = make([]HTTPMiddleware, 0, len(s5f4mw))
	}
	p7this.s5f4middleware = append(p7this.s5f4middleware, s5f4mw...)
}

// Get 装饰器模式 包装addroute
func (p7this *HTTPHandler) Get(path string, f4h HTTPHandlerFunc, s5f4mw ...HTTPMiddleware) {
	p7this.router.addRoute(http.MethodGet, path, f4h, s5f4mw...)
}

func (p7this *HTTPHandler) Post(path string, f4h HTTPHandlerFunc, s5f4mw ...HTTPMiddleware) {
	p7this.router.addRoute(http.MethodPost, path, f4h, s5f4mw...)
}
