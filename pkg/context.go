package pkg

import (
	"encoding/json"
	"io"
	"net/http"
)

// HTTPContext 封装 Handler.ServeHTTP 方法的两个参数
// 这里是个简单的实现，想要完善一点的，可以实现 context.Context 接口
type HTTPContext struct {
	P1resW http.ResponseWriter
	P1req  *http.Request
}

func NewHTTPContext(p1resW http.ResponseWriter, p1req *http.Request) *HTTPContext {
	return &HTTPContext{
		P1resW: p1resW,
		P1req:  p1req,
	}
}

// ReadJson 读取数据转换为 json
func (p1c *HTTPContext) ReadJson(data interface{}) error {
	reqBody, err := io.ReadAll(p1c.P1req.Body)
	if nil != err {
		return err
	}
	return json.Unmarshal(reqBody, data)
}

// WriteJson 写入 json 数据
func (p1c *HTTPContext) WriteJson(status int, data interface{}) error {
	p1c.P1resW.WriteHeader(status)
	resJson, err := json.Marshal(data)
	if nil != err {
		return err
	}
	_, err = p1c.P1resW.Write(resJson)
	if nil != err {
		return err
	}
	return nil
}
