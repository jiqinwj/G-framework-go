package pkg

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

type HTTPShutdown struct {
	// 服务是否关闭，0==运行，1==关闭
	isClose int32
	// 正在处理的请求数量
	reqInHandle int32
	// 服务关闭时负责接收请求处理完毕的信号
	canShutdown chan struct{}
}

func NewHTTPShutdown() *HTTPShutdown {
	return &HTTPShutdown{
		canShutdown: make(chan struct{}),
	}
}

// ReqInHandleCountBuilder 处理请求（关闭时拒绝新请求，记录正在处理的请求数量）
func (p1s *HTTPShutdown) ReqInHandleCountBuilder(next MiddlewareFunc) MiddlewareFunc {
	return func(p1c *HTTPContext) {
		isClose := atomic.LoadInt32(&p1s.isClose)

		//如果服务已经关闭，就拒绝新的请求
		if isClose > 0 {
			p1c.P1resW.WriteHeader(http.StatusNotFound)
			_, _ = p1c.P1resW.Write([]byte(fmt.Sprintf("service is closing")))
			return
		}

		atomic.AddInt32(&p1s.reqInHandle, 1)
		next(p1c)
		reqInHandle := atomic.AddInt32(&p1s.reqInHandle, -1)

		// 如果服务已经关闭，正在处理的请求数量为0，就可以发送请求处理完毕的信号
		if isClose > 0 && reqInHandle == 0 {
			p1s.canShutdown <- struct{}{}
		}

	}
}

// RejectNewRequestAndWaiting 拒绝新的请求，等待正在处理的请求处理完成
func (p1s *HTTPShutdown) RejectNewRequestAndWaiting(c context.Context) error {
	atomic.AddInt32(&p1s.isClose, 1)

	//关闭之前，就已经处理完了请求，就不需要等了，直接返回。
	if atomic.LoadInt32(&p1s.reqInHandle) == 0 {
		return nil
	}

	select {
	case <-p1s.canShutdown:
		fmt.Println("RejectNewRequestAndWaiting,HTTPShutdown canShutdown")
	case <-c.Done():
		return errors.New("RejectNewRequestAndWaiting,Context Done")
	}
	return nil
}

func WaitForShutdown(arr1hookFunc ...HookFunc) {
	chansignal := make(chan os.Signal, 1)
	// os.Interrupt = syscall.SIGINT，就是Ctrl+C
	signal.Notify(chansignal, os.Interrupt)

	select {
	case t1signal := <-chansignal:
		fmt.Printf("WaitForShutdown,get signal %s,service will shutdown\r\n", t1signal)
		// 十分钟都没关闭，就强行关闭
		time.AfterFunc(time.Minute*10, func() {
			fmt.Printf("WaitForShutdown,shutdown timeout,service will shutdown immediately\r\n")
			os.Exit(1)
		})
		// 依次执行 hook
		for _, t1hookFunc := range arr1hookFunc {
			c, cancel := context.WithTimeout(context.Background(), time.Second*30)
			err := t1hookFunc(c)
			if err != nil {
				fmt.Printf("WaitForShutdown,hook failed,err %+v\r\n", err)
			}
			cancel()
		}
		os.Exit(0)
	}

}
