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
		p7ctx.Reset()
		p7this.ctxPool.Put(p7ctx)
	}()
	p7ctx.I9writer = i9w
	p7ctx.P7request = p7r

	// 倒过来组装，先组装的在里层，里层的后执行
	// 最里层应该是找路由然后执行业务代码
	//t4chain := p7this.do

}

func (p7this *HTTPHandler) doServerHTTP(p7ctx *HTTPContext) {
	// 如果服务已经关闭了就直接返回
	if !p7this.isRunning {
		p7ctx.I9writer.WriteHeader(http.StatusInternalServerError)
		_, _ = p7ctx.I9writer.Write([]byte("服务已关闭"))
		return
	}
	//p7ri := p7this.findRoute(p7ctx.P7request.Method, p7ctx.P7request.URL.Path)
	// 如果找不到对应的路由结点或者路由结点上没有处理方法就返回 404
	//if nil== p7ri || nil == p7ri

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

// 路由数据
type RouteData struct {
	Method   string
	Path     string
	F4handle HTTPHandlerFunc
}

// 添加一组路由
func (p7this *HTTPHandler) Group(path string, s5f4mw []HTTPMiddleware, s5routeData []RouteData) {

	for _, rd := range s5routeData {
		t4path := path
		if "/" != rd.Path {
			t4path = path + rd.Path
		}
		p7this.addRoute(rd.Method, t4path, rd.F4handle, s5f4mw...)
	}
}
