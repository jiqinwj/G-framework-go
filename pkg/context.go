package pkg

import (
	"fmt"
	"net/http"
)

// HTTPContext 自定义请求上下文（注意和 context.Context 的概念区分开）
type HTTPContext struct {
	// ServeHTTP 的 http.ResponseWriter
	I9writer http.ResponseWriter
	// ServeHTTP 的 *http.Request
	P7request *http.Request

	// 命中的路由结点
	p7routingNode *routingNode
	// 提取到的路径参数
	M3pathParam map[string]string

	// RespData 暂存请求数据
	// 因为 http.Request.Body 是流，只能读一次（和 linux c 的 recvFrom() 类似）
	// 如果等到应用层再调用，那么在中间件里面记录请求日志或者进行预处理就无法实现
	// 这里的方案是，在所有的处理流程开始前就读取然后存下来，如果有需要可以再造一个流放回去
	ReqBody []byte
	// RespData 暂存响应 http status code
	RespStatusCode int
	// RespData 暂存响应数据
	// 因为 http.ResponseWriter.Write 是流，只能写一次（和 linux c 的 write() 类似）
	// 如果在应用层调用了，那么在中间件里面记录响应日志或者追加数据就无法实现
	// 这里的方案是，等到所有的处理流程都结束了，再调用 http.ResponseWriter.Write
	RespData []byte
}

func NewHTTPContext() *HTTPContext {
	return &HTTPContext{
		// 一般情况下，路径参数都是 1 个
		M3pathParam: make(map[string]string, 1),
	}
}

// Reset 复用时用新的数据重置
func (p7this *HTTPContext) Reset() {
	p7this.I9writer = nil
	p7this.P7request = nil
	p7this.p7routingNode = nil
	p7this.M3pathParam = make(map[string]string, 1)
	p7this.ReqBody = nil
	p7this.RespStatusCode = 0
	p7this.RespData = nil
}

// GetRoutingInfo 获取命中的路由结点的信息
func (this HTTPContext) GetRoutingInfo() string {
	return fmt.Sprintf("nodeType:%d;routing path:%s;", this.p7routingNode.nodeType, this.p7routingNode.path)
}
