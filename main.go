package main

import (
	"fmt"
	"ji-framework-go/middleware"
	"ji-framework-go/pkg"
	"ji-framework-go/shutdown"
	"net/http"
	"time"
)

func main() {
	p7os := makeOpenService()
	p7as := makeAdminService()
	p7sm := pkg.NewServiceManager(
		[]*pkg.HTTPService{p7os, p7as},
		pkg.SetShutdownTimeOutOption(20*time.Second),
		pkg.SetShutdownWaitTime(10*time.Second),
		pkg.SetShutdownCallbackTimeOut(5*time.Second),
	)
	p7sm.Start()
}

func makeOpenService() *pkg.HTTPService {
	p7h := pkg.NewHTTPHandler()

	p7h.AddMiddleware(
		middleware.RecoveryMiddleware(),
		middleware.ReqBodyMiddleware(),
		middleware.LogMiddleware(),
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

	p7h.Get("/order", f4handler)
	p7h.Get("/order/list/:size/:page", f4handler)
	p7h.Get("/order/:id/detail", f4handler)
	p7h.Post("/order/create", f4handler)
	p7h.Post("/order/:id/delete", f4handler)

	p7s := pkg.NewHTTPService("9510", "127.0.0.1:9510", p7h)

	p7s.AddShutdownCallback(
		shutdown.CacheShutdownCallback,
		shutdown.CountShutdownCallback,
	)

	return p7s
}

func makeAdminService() *pkg.HTTPService {
	p7h := pkg.NewHTTPHandler()

	p7h.AddMiddleware(
		middleware.RecoveryMiddleware(),
		middleware.ReqBodyMiddleware(),
		middleware.LogMiddleware(),
	)

	f4handler := func(p7ctx *pkg.HTTPContext) {
		routingInfo := p7ctx.GetRoutingInfo()
		pathParam := "pathParam:"
		for key, val := range p7ctx.M3pathParam {
			pathParam += fmt.Sprintf("%s=%s;", key, val)
		}
		p7ctx.RespData = append(p7ctx.RespData, []byte(routingInfo+pathParam)...)
	}

	p7h.Group(
		"/admin",
		[]pkg.HTTPMiddleware{middleware.TestMiddleware("admin")},
		[]pkg.RouteData{
			{http.MethodGet, "/", f4handler},
			{http.MethodGet, "/list/:size/:page", f4handler},
			{http.MethodGet, "/:id/detail", f4handler},
			{http.MethodPost, "/create", f4handler},
			{http.MethodPost, "/:id/delete", f4handler},
		},
	)

	return pkg.NewHTTPService("9511", "127.0.0.1:9511", p7h)
}
