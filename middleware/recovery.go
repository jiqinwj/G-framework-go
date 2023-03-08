package middleware

import "ji-framework-go/pkg"

func RecoverMiddleware() pkg.HTTPMiddleware {
	return func(next pkg.HTTPHandlerFunc) pkg.HTTPHandlerFunc {
		return func(p7ctx *pkg.HTTPContext) {
			defer func() {
				if err := recover(); err != nil {
					p7ctx.RespData = append(p7ctx.RespData, []byte("recover from panic\r\n")...)
				}
			}()
		}
	}
}
