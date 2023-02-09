package pkg

import "fmt"

// MiddlewareFunc 中间件处理类
type MiddlewareFunc func(c *HTTPContext)

// MiddlewareBuilder 中间件建造器
// 主要实现思路 就是责任链模式。一个套一个 最后返回一个 MiddlewareFunc
// 在返回的方法内部会调用传入的nex MiidlewareFunc
type MiddlewareBuilder func(next MiddlewareFunc) MiddlewareFunc

// TestMiddlerwareBuilder 测试调用顺序
func TestMiddlerwareBuilder(next MiddlewareFunc) MiddlewareFunc {
	return func(c *HTTPContext) {
		fmt.Printf("request before test middleware.\n")
		next(c)
		fmt.Printf("request after test middleware.\n")
	}
}
