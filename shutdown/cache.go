package shutdown

import (
	"context"
	"log"
	"time"
)

// 持久化缓存内容
func CacheShutdownCallback(ctx context.Context) {
	c4signal := make(chan struct{}, 1)
	go func() {
		log.Println("持久化缓存内容中。。。")
		time.Sleep(1 * time.Second)
		c4signal <- struct{}{}
	}()
	select {
	case <-c4signal:
		log.Println("持久化缓存内容成。")
	case <-ctx.Done():
		log.Println("持久化缓存内容超时。")
	}

}
