package main

import (
	"fmt"
	"ji-framework-go/middleware"
	"ji-framework-go/pkg"
)

func main() {
	p7os := makeOpenService()
}

func makeOpenService() *pkg.HTTPService {
	p7h := pkg.NewHTTPHandler()

	p7h.AddMiddleware(
		middleware.RecoverMiddleware(),
	)

	f4handler := func(p7ctx *pkg.HTTPContext) {
		routingInfo := p7ctx.GetRoutingInfo()
		pathParam := "pathParam:"
		for key, val := range p7ctx.M3pathParam {
			pathParam += fmt.Sprintf("%s=%s;", key, val)
		}
		p7ctx.RespData = append(p7ctx.RespData, []byte(routingInfo+pathParam)...)
	}

	p7h.Get("/", f4handler)

	p7h.Get("/hello", f4handler)
	p7h.Get("/hello/world", f4handler, middleware.TestMiddleware("/hello"), middleware.TestMiddleware("/world"))
	p7h.Get("/hello/*", f4handler, middleware.TestMiddleware("/hello/*"))

	//启动服务
	p7s := pkg.NewHTTPService("9510", "127.0.0.1:9510", p7h)

}
