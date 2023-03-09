package middleware

import (
	"bytes"
	"io"
	"ji-framework-go/pkg"
)

// ReqBodyMiddleware 提取请求参数
func ReqBodyMiddleware() pkg.HTTPMiddleware {
	return func(next pkg.HTTPHandleFunc) pkg.HTTPHandleFunc {
		return func(p7ctx *pkg.HTTPContext) {
			// 处理请求参数
			var err error
			p7ctx.ReqBody, err = io.ReadAll(p7ctx.P7request.Body)
			if nil != err {
				return
			}
			p7ctx.P7request.Body = io.NopCloser(bytes.NewBuffer(p7ctx.ReqBody))

			next(p7ctx)
		}
	}
}
