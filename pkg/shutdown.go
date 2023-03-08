package pkg

import "context"

// 服务关闭时需要执行的回调方法的定义
type ShutdownCallback func(ctx context.Context)

// 服务关闭时需要执行的回调方法的接口定义
type ShutdownCallbackInterface interface {
	// AddShutdownCallback 添加服务关闭时需要执行的回调方法
	AddShutdownCallback(s5 ...ShutdownCallback)
}
