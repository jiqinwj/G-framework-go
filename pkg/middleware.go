package pkg

import (
	"fmt"
	"time"
)

// MiddlewareFunc 中间件处理方法。
// 这里和 HTTPHandlerFunc 保持一致。要不然最后一环，调用没法串起来。
type MiddlewareFunc func(c *HTTPContext)

// MiddlewareBuilder 中间件建造器。
// 实现思路就是链式套娃。返回一个 MiddlewareFunc 方法。
// 在返回的方法内部会调用传入的 next MiddlewareFunc 方法。
type MiddlewareBuilder func(next MiddlewareFunc) MiddlewareFunc

// Test1MiddlewareBuilder 测试调用顺序
func TestMiddlewareBuilder(next MiddlewareFunc) MiddlewareFunc {
	return func(c *HTTPContext) {
		fmt.Printf("request before test middleware.\n")
		next(c)
		fmt.Printf("request after test middleware.\n")
	}
}

// TimeCostMiddlewareBuilder 算一下耗时
func TimeCostMiddlewareBuilder(next MiddlewareFunc) MiddlewareFunc {
	return func(c *HTTPContext) {
		startUN := time.Now().UnixNano()
		next(c)
		endUN := time.Now().UnixNano()
		fmt.Printf("request time cost: %d unix nano.\n", startUN-endUN)
	}
}
