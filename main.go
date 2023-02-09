package main

import (
	jiframework "ji-framework-go/pkg"
	"net/http"
)

type ApiJson struct {
	JsonInt    int    `json:"json_int"`
	JsonString string `json:"json_string"`
	JsonText   string `json:"json_text"`
}

func main() {

	//创建一个service 启动服务类
	n2hservice := jiframework.NewHTTPService(
		"http-service",
		jiframework.TestMiddlerwareBuilder,
	)
	//注册
	registerApi(n2hservice)
	//启动服务了
	n2hservice.Start("127.0.0.1", "9504")
}

// registerApi 注册路由和处理方法
func registerApi(n2hservice jiframework.Service) {
	n2hservice.RegisteRoute(http.MethodGet, "/api/test", func(n2c *jiframework.HTTPContext) {
		n2c.N1resW.WriteHeader(http.StatusOK)
		_, _ = n2c.N1resW.Write([]byte("reponse,http.MethodGet,/api/test"))
	})
}
