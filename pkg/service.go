package pkg

import (
	"fmt"
	"net/http"
)

// Service 服务接口
type Service interface {
	//Start 服务开始了
	Start(addr string, port string) error
	//服务肯定要注册路由啦
	HTTPRoute
}

// HTTPService Http 服务启动结构体
type HTTPService struct {
	//Name 服务的花名
	Name string
	//HttpHandler 服务肯定需要一个处理器啊。好比汽车需要引擎启动啊
	handler HTTPHandler
	//中间件 🧅洋葱模型哦 入口方法
	entrance MiddlewareFunc
}

// NewHTTPService 创建一个 Service 接口的类。指定下服务的花名和中间件组合
func NewHTTPService(name string, arr1Builder ...MiddlewareBuilder) Service {
	//基于这个路由树
	//要实例化 一个handler 一个路由树
	var n2h HTTPHandler = NewHTTPHandlerTree()

	//在使用中间件的时候，需要对请求处理的入口方法进行封装
	var hf MiddlewareFunc = n2h.HandlerHTTP
	//反过来遍历中间件建造器数组，像洋葱一样，数组最前面的对应最外层
	// 套娃完成后，请求处理的入口方法在里面，表示请求通过层层中间件后进入业务逻辑
	for i := len(arr1Builder) - 1; i > -1; i-- {
		var mf MiddlewareBuilder = arr1Builder[i]
		hf = mf(hf)
	}

	return &HTTPService{
		Name:     name,
		handler:  n2h,
		entrance: hf,
	}

}

// ServerHTTP Handler.ServerHTTP 把HTTPService 结构体变成 src/net/http/service.go 里handler 接口的实现列
// 在调用 http.ListenAndServe(addr string,handler Handler) 的时候，会把HTTPService 的实列作为handler 参数传入进去
// http.ListenAndServe() 会创建一个 src/net/http/server.go 里 Server 结构体的实列，保存 hanler 参数
// 然后 http.ListenAndServe() 会调用 Net.listen(network,address string) 启动TCP 服务
// net.Listen() 返回一个 net.Listener 接口的实例，net.Listener 实例通过 Accept() 方法获取 TCP 连接。
// 获取到 TCP 连接之后，经过一系列的操作，最后会有这么一行代码 serverHandler{c.server}.ServeHTTP(w, w.req)。
// 这行代码会把一开始传进去的Handler 接口的实列（HTTPService的实列）取出来。然后调用ServerHttp 方法
func (n2s *HTTPService) ServeHTTP(n2resw http.ResponseWriter, n2req *http.Request) {
	p1c := NewHTTPContext(n2resw, n2req)

	//不使用中间件时，直接调用 HTTPHandler 的实列处理请求
	//n2s.handler.HandlerHTTP(p1c)

	//使用中间件后，这里就要改成调用中间件入口
	n2s.entrance(p1c)
}

// Start Service.start
func (n2s *HTTPService) Start(addr string, port string) error {
	fmt.Printf("HTTPService %s start at %s...\n", n2s.Name, addr+":"+port)
	return http.ListenAndServe(addr+":"+port, n2s)
}

// RegisteRoute Service.HTTPRute.RegisteRoute
func (n2s *HTTPService) RegisteRoute(method string, pattern string, hhFunc HTTPHandlerFunc) error {
	return n2s.handler.RegisteRoute(method, pattern, hhFunc)
}
