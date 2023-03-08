package pkg

import (
	"context"
	"net/http"
)

// HTTPServiceInterface 核心服务的接口定义
type HTTPServiceInterface interface {
	// Start 启动服务
	Start(addr string) error
	// Stop 停止服务
	// ShutDown 服务关闭
	ShutDown(ctx context.Context) error
	ShutdownCallbackInterface
}

// 核心服务
type HTTPService struct {
	// name 服务名
	name string
	// 核心处理逻辑
	p7server *http.Server
	// 核心处理逻辑
	p7handler *HTTPHandler
	// 服务是否正在运行
	isRunning bool
	// 关闭服务时，需要执行的回调方法
	s5f4shutdownCallback []ShutdownCallback
}

func NewHTTPService(name string, addr string, p7h *HTTPHandler) *HTTPService {
	return &HTTPService{
		name: name,
		p7server: &http.Server{
			Addr:    addr,
			Handler: p7h,
		},
		p7handler: p7h,
		isRunning: true,
	}
}
