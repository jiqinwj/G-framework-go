package pkg

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type HookFunc func(c context.Context) error

// NotifyShutdownToGateWay 通知网关，服务要下线了
func NotifyShutdownToGateWay(p1c context.Context) error {
	fmt.Printf("notify gateway server will shutdown")
	time.Sleep(time.Second * 2)
	return nil
}

// ServiceShutDownBuilder 关闭服务
func ServiceShutDownBuilder(arr1p1s ...Service) HookFunc {
	return func(c context.Context) error {
		syncWg := sync.WaitGroup{}
		chanSyncWg := make(chan struct{})

		syncWg.Add(len(arr1p1s))
		for _, t1p1s := range arr1p1s {
			go func(p1s Service) {
				err := p1s.Shutdown(c)
				if err != nil {
					fmt.Printf("service shutdown err: %+v\r\n", err)
				}
				time.Sleep(time.Second)
				syncWg.Done()
			}(t1p1s)
		}

		go func() {
			syncWg.Wait()
			chanSyncWg <- struct{}{}
		}()

		select {
		case <-chanSyncWg:
			fmt.Printf("SerivceShutDown,close all servers\r\n")
			return nil
		case <-c.Done():
			fmt.Printf("ServiceShutdown,Context Done\r\n")
			return errors.New("ServiceShutdown,Context Done")
		}
	}
}
