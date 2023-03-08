package pkg

// 中间件的处理方法的定义
type HTTPMiddleware func(next HTTPHandlerFunc) HTTPHandlerFunc

// MiddlewareInterface 中间件接口定义
type MiddlewareInterface interface {
	// AddMiddleware 添加中间件
	AddMiddleware(s5mw ...HTTPMiddleware)
}
