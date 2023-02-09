package pkg

import (
	"fmt"
	"net/http"
)

// Service 服务接口
type Service interface {
	// Start 服务启动
	Start(addr string, port string) error
	// 服务接口需要能够注册路由
	HTTPRoute
}

// HTTPService HTTP 服务
type HTTPService struct {
	// Name 服务的名字
	Name string
	// HTTPHandler 服务需要一个请求处理接口的实例
	handler HTTPHandler
	// 中间件入口方法
	entrance MiddlewareFunc
}

// NewHTTPSrevice 创建一个 Service 接口的实例，指定服务的名字和中间件组
func NewHTTPSrevice(name string, arr1Builder ...MiddlewareBuilder) Service {
	// 这里实例化一个 HTTPHandler 接口的实例
	// var p1h HTTPHandler = NewHTTPHandlerMap()
	var p1h HTTPHandler = NewHTTPHandlerTree()

	// 在使用中间件的时候，需要对请求处理的入口方法进行包装。
	var hf MiddlewareFunc = p1h.HandlerHTTP
	// 反过来遍历中间件建造器数组。像洋葱一样，数组最前面的对应最外层。
	// 套娃完成后，请求处理的入口方法应该在最里面。表示请求通过层层中间件后进入业务逻辑。
	for i := len(arr1Builder) - 1; i > -1; i-- {
		var mf MiddlewareBuilder = arr1Builder[i]
		hf = mf(hf)
	}

	return &HTTPService{
		Name:     name,
		handler:  p1h,
		entrance: hf,
	}
}

// ServeHTTP Handler.ServeHTTP，把 HTTPService 结构体变成 src/net/http/service.go 里 Handler 接口的实例。
// 在调用 http.ListenAndServe(addr string, handler Handler) 的时候，会把 HTTPService 的实例作为 handler 参数传进去。
// http.ListenAndServe() 会创建一个 src/net/http/server.go 里 Server 结构体的实例，保存 handler 参数。
// 然后 http.ListenAndServe() 调用 net.Listen(network, address string)，启动 TCP 服务。
// net.Listen() 返回一个 net.Listener 接口的实例，net.Listener 实例通过 Accept() 方法获取 TCP 连接。
// 获取到 TCP 连接之后，经过一系列的操作，最后会有这么一行代码 serverHandler{c.server}.ServeHTTP(w, w.req)。
// 这行代码会把一开始传进去的 Handler 接口的实例（HTTPService 的实例）取出来，然后调用 ServeHTTP 方法。
func (p1s *HTTPService) ServeHTTP(p1resW http.ResponseWriter, p1req *http.Request) {
	p1c := NewHTTPContext(p1resW, p1req)

	// 不使用中间件时，直接调用 HTTPHandler 的实例处理请求
	// p1s.handler.HandlerHTTP(p1c)

	// 使用中间件后，这里就要改成调用中间件入口
	p1s.entrance(p1c)
}

// Start Service.Start
func (p1s *HTTPService) Start(addr string, port string) error {
	fmt.Printf("HTTPService %s start at %s...\n", p1s.Name, addr+":"+port)
	return http.ListenAndServe(addr+":"+port, p1s)
}

// RegisteRoute Service.HTTPRoute.RegisteRoute
func (p1s *HTTPService) RegisteRoute(method string, pattern string, hhFunc HTTPHandlerFunc) error {
	return p1s.handler.RegisteRoute(method, pattern, hhFunc)
}
