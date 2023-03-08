package middleware

import "ji-framework-go/pkg"

func TestMiddleware(code string) pkg.HTTPMiddleware {
	code = "[" + code + "]"
	return func(next pkg.HTTPHandlerFunc) pkg.HTTPHandlerFunc {
		return func(p7ctx *pkg.HTTPContext) {
			p7ctx.RespData = append(p7ctx.RespData, []byte(code)...)
			next(p7ctx)
		}
	}
}
