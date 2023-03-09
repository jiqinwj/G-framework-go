package pkg

import "context"

// ShutdownCallback 服务关闭时需要执行的回调方法的定义
type ShutdownCallback func(context.Context)

// ShutdownCallbackInterface 服务关闭时需要执行的回调方法的接口定义
type ShutdownCallbackInterface interface {
	// AddShutdownCallback 添加服务关闭时需要执行的回调方法
	AddShutdownCallback(s5f4cb ...ShutdownCallback)
}
