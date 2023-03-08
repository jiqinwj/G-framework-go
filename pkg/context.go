package pkg

import (
	"fmt"
	"net/http"
)

// 自定义请求上下文 (注意和context.Context 的概念区分开)
type HTTPContext struct {
	//ServerHTTP 的 http.ResponseWriter
	I9writer http.ResponseWriter
	// ServerHTTP 的 *http.Request
	P7request *http.Request

	// 命中的路由结点
	p7routingNode *routingNode

	// 提取到的路径参数
	M3pathParam map[string]string

	// 暂存请求数据
	//因为 http.Request.Body 是流，只能读一次 (和linux c 的recvForm() 类似)
	//如果等到应用层再调用，那么在中间件里面记录的请求日志或者进行预处理就无法实现
	//这里的方案是，在所有的处理流程开始前就读取然后存下来，如果有需求可以再造一个流放回去
	ReqBody []byte
	// 暂存响应 http status code
	RespStatusCode int

	// 因为 http.ResponseWriter.Write 是流，只能写一次 （其实内部 就是linux c 的write() 类似）
	// 如果在应用层调用了，那么在中间件里面记录的响应日志和追加数据就无法实现
	// 这里的方案是，等到所有的处理流程都结束了，再调用 http.ResponseWriter.Writer
	RespData []byte
}

func NewHTTPContext() *HTTPContext {
	return &HTTPContext{
		// 一般情况下，路径参数就一个
		M3pathParam: make(map[string]string, 1),
	}
}

// GetRoutingInfo  获取命中的路由结点的信息
func (this HTTPContext) GetRoutingInfo() string {
	return fmt.Sprintf("nodeType:%d;routing path:%s;", this.p7routingNode.nodeType, this.p7routingNode.path)
}

// 复用时用新的数据重置
func (p7this *HTTPContext) Reset() {
	p7this.I9writer = nil
	p7this.P7request = nil
	p7this.p7routingNode = nil
	p7this.M3pathParam = make(map[string]string, 1)
	p7this.ReqBody = nil
	p7this.RespStatusCode = 0
	p7this.RespData = nil
}
