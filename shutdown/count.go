package shutdown

import (
	"context"
	"log"
	"time"
)

// CountShutdownCallback 上报统计数据
func CountShutdownCallback(ctx context.Context) {
	c4signal := make(chan struct{}, 1)
	go func() {
		log.Println("上报统计数据。。。")
		time.Sleep(1 * time.Second)
		c4signal <- struct{}{}
	}()
	select {
	case <-c4signal:
		log.Println("上报统计数据完成。")
	case <-ctx.Done():
		log.Println("上报统计数据超时。")
	}
}
