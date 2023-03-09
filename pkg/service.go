package pkg

import (
	"context"
	"log"
	"net/http"
)

// HTTPServiceInterface 核心服务的接口定义
type HTTPServiceInterface interface {
	// Start 启动服务
	Start(addr string) error
	// Stop 停止服务
	Stop()
	// ShutDown 服务关闭
	ShutDown(ctx context.Context) error
	ShutdownCallbackInterface
}

// HTTPService 核心服务
type HTTPService struct {
	// name 服务名
	name string
	// p7server 核心处理逻辑
	p7server *http.Server
	// p7handler 核心处理逻辑
	p7handler *HTTPHandler
	// isRunning 服务是否正在运行
	isRunning bool
	// s5f4shutdownCallback 关闭服务时，需要执行的回调方法
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

func (p7this *HTTPService) Start() error {
	log.Printf("服务 %s 启动，监听 %s 端口。\r\n", p7this.name, p7this.p7server.Addr)
	p7this.p7handler.router.middlewareCache()
	return p7this.p7server.ListenAndServe()
}

func (p7this *HTTPService) Stop() {
	log.Printf("服务 %s 停止服务\r\n", p7this.name)
	p7this.isRunning = false
	p7this.p7handler.isRunning = false
}

func (p7this *HTTPService) ShutDown(ctx context.Context) error {
	log.Printf("服务 %s 关闭\r\n", p7this.name)
	return p7this.p7server.Shutdown(ctx)
}

func (p7this *HTTPService) AddShutdownCallback(s5f4cb ...ShutdownCallback) {
	if nil == p7this.s5f4shutdownCallback {
		p7this.s5f4shutdownCallback = make([]ShutdownCallback, 0, len(s5f4cb))
	}
	p7this.s5f4shutdownCallback = append(p7this.s5f4shutdownCallback, s5f4cb...)
}
