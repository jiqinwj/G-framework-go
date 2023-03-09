package middleware

import (
	"bytes"
	"io"
	"ji-framework-go/pkg"
)

/*
#note NopCloser 的原理很简单，就是将一个不带 Close 的 Reader 封装成 ReadCloser
@link https://blog.csdn.net/DisMisPres/article/details/123380773
*/
func ReqBodyMiddleware() pkg.HTTPMiddleware {
	return func(next pkg.HTTPHandlerFunc) pkg.HTTPHandlerFunc {
		return func(p7ctx *pkg.HTTPContext) {
			//处理请求参数
			var err error
			p7ctx.ReqBody, err = io.ReadAll(p7ctx.P7request.Body)
			if nil != err {
				return
			}
			//注意这时 p7ctx.P7request.Body 已经读完了，需要重新将读出来的值给放回去，之后的处理就依然可以使用 c.Request.Body 了
			p7ctx.P7request.Body = io.NopCloser(bytes.NewBuffer(p7ctx.ReqBody))
			next(p7ctx)
		}
	}
}
