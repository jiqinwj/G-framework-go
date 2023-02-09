package pkg

import (
	"encoding/json"
	"io"
	"net/http"
)

type HTTPContext struct {
	N1resW http.ResponseWriter
	N1req  *http.Request
}

func NewHTTPContext(N1resW http.ResponseWriter, N1req *http.Request) *HTTPContext {
	return &HTTPContext{
		N1resW: N1resW,
		N1req:  N1req,
	}
}

// ReadJson 读取数据转换为Json
func (n2c *HTTPContext) ReadJson(data interface{}) error {
	reqBody, err := io.ReadAll(n2c.N1req.Body)
	if nil != err {
		return err
	}
	return json.Unmarshal(reqBody, data)
}

// WriteJson 写入json 数据
func (n2c *HTTPContext) WriteJson(status int, data interface{}) error {
	n2c.N1resW.WriteHeader(status)
	resJson, err := json.Marshal(data)
	if nil != err {
		return err
	}
	_, err = n2c.N1resW.Write(resJson)
	if nil != err {
		return err
	}
	return nil
}
