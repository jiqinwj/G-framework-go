package middleware

import "ji-framework-go/pkg"

func RecoveryMiddleware() pkg.HTTPMiddleware {
	return func(next pkg.HTTPHandleFunc) pkg.HTTPHandleFunc {
		return func(p7ctx *pkg.HTTPContext) {
			defer func() {
				if err := recover(); err != nil {
					p7ctx.RespData = append(p7ctx.RespData, []byte("recover from panic\r\n")...)
				}
			}()
			next(p7ctx)
		}
	}
}
