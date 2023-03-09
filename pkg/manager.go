package pkg

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

// option 设计模式
type ServiceManagerOption func(*ServiceManager)

type ServiceManager struct {
	// 管理的服务列表
	s5p7HTTPService []*HTTPService
	// 服务关闭的总超时时间，超时强制关闭
	shutdownTimeOut time.Duration
	// 服务关闭时，等待正在处理的请求的时间
	// 等待结束后，开始执行服务关闭时需要执行的回调方法
	shutdownWaitTime time.Duration
	// 服务关闭时需要执行的回调方法的超时时间
	shutdownCallbackTimeOut time.Duration
}

func NewServiceManager(s5p7hs []*HTTPService, s5f4option ...ServiceManagerOption) *ServiceManager {
	p7sm := &ServiceManager{
		s5p7HTTPService:         s5p7hs,
		shutdownTimeOut:         10 * time.Second,
		shutdownWaitTime:        3 * time.Second,
		shutdownCallbackTimeOut: 3 * time.Second,
	}

	// 依次执行option
	for _, f4o := range s5f4option {
		f4o(p7sm)
	}

	return p7sm
}

func (p7this *ServiceManager) Start() {
	//启动服务
	log.Println("服务启动中.....")
	for _, p7s := range p7this.s5p7HTTPService {
		t4p7s := p7s
		go func() {
			if err := t4p7s.Start(); nil != err {
				if http.ErrServerClosed == err {
					log.Printf("子服务 %s 已关闭\n", t4p7s.name)
				} else {
					log.Printf("子服务 %s 异常退出，err:%s\r\n", t4p7s.name, err)
				}
			}
		}()
	}
	log.Println("服务启动完成。。。")
	// 监听 ctrl+c 信号
	c4signal := make(chan os.Signal, 2)
	signal.Notify(c4signal, os.Interrupt)
	select {
	case <-c4signal:
		log.Printf("接收到关闭信号，开始关闭服务，限制 %d 秒内完成。。。。", p7this.shutdownWaitTime/time.Second)

		//再次监听 ctrl+C 信号
		go func() {
			select {
			case <-c4signal:
				log.Println("再次接收到关闭信号，服务直接退出。")
				os.Exit(1)
			}
		}()
		time.AfterFunc(p7this.shutdownTimeOut, func() {
			log.Println("优雅关闭超时，服务直接退出")
			os.Exit(1)
		})
		p7this.Shutdown()
	}
}

func (p7this *ServiceManager) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), p7this.shutdownTimeOut)
	defer cancel()

	log.Println("停止接收新请求。。")
	for _, p7hs := range p7this.s5p7HTTPService {
		p7hs.Stop()
	}

	log.Printf("等待正在执行的请求结束，等待%d 秒。。。。", p7this.shutdownWaitTime/time.Second)
	time.Sleep(p7this.shutdownWaitTime)

	log.Println("开始关闭子服务。。。。")
	wg := sync.WaitGroup{}
	for _, p7hs := range p7this.s5p7HTTPService {
		log.Printf("关闭子服务 %s....", p7hs.name)
		t4p7s := p7hs
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = t4p7s.ShutDown(ctx)
		}()
	}
	wg.Wait()

	log.Println("开始执行子服务的关闭回调。。。。")
	for _, p7hs := range p7this.s5p7HTTPService {
		log.Printf("执行子服务 %s 的关闭回调，限制 %d 秒内完成。。。。", p7hs.name, p7this.shutdownCallbackTimeOut/time.Second)
		for _, f4cb := range p7hs.s5f4shutdownCallback {
			t4f4cb := f4cb
			wg.Add(1)
			go func() {
				defer wg.Done()
				t4ctx, t4cancel := context.WithTimeout(context.Background(), p7this.shutdownCallbackTimeOut)
				defer t4cancel()
				t4f4cb(t4ctx)
			}()
		}
	}
	wg.Wait()

	log.Println("服务关闭完成。。")

}

func SetShutdownTimeoutOption(t time.Duration) ServiceManagerOption {
	return func(p7sm *ServiceManager) {
		p7sm.shutdownTimeOut = t
	}
}

func SetShutdownWaitTime(t time.Duration) ServiceManagerOption {
	return func(p7sm *ServiceManager) {
		p7sm.shutdownWaitTime = t
	}
}

func SetShutdownCallbackTimeOut(t time.Duration) ServiceManagerOption {
	return func(p7sm *ServiceManager) {
		p7sm.shutdownCallbackTimeOut = t
	}
}
