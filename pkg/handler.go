package pkg

// HTTPHandler 请求处理接口
type HTTPHandler interface {
	//HandlerHTTP 请求处理的入口方法
	// 所有请求必须都先进入这个方法进行路由匹配，然后调用路由对应的处理方法
	HandlerHTTP(c *HTTPContext)
	// 请求处理接口肯定需要注册路由
	HTTPRoute
}

// HTTPRoute 路由接口
type HTTPRoute interface {
	//RegisteRoute 注册路由，method HTTP 方法：pattern 路由
	RegisteRoute(method string, pattern string, hhFunc HTTPHandlerFunc) error
}

// HTTPHandlerFunc 路由对应的处理方法
type HTTPHandlerFunc func(c *HTTPContext)
