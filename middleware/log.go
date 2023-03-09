package middleware

import (
	"fmt"
	"ji-framework-go/pkg"
)

func LogMiddleware() pkg.HTTPMiddleware {
	return func(next pkg.HTTPHandlerFunc) pkg.HTTPHandlerFunc {
		return func(p7ctx *pkg.HTTPContext) {
			fmt.Printf("request path:%s\r\n", p7ctx.P7request.URL.Path)
			fmt.Println("ReqBody:", string(p7ctx.ReqBody))
			next(p7ctx)
			fmt.Println("RespData:", string(p7ctx.RespData))
		}
	}
}
