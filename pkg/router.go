package pkg

import (
	"strings"
)

var (
	StrPathCannotBeEmpty      = "路由不能为空"
	StrPartCannotBeEmpty      = "路由不能存在连续的 `\\`"
	StrPathMustStartWithSlash = "路由必须以 `\\` 开始"
	StrPathCannotEndWithSlash = "路由不能以 `\\` 结尾"

	StrRootNodeExist = "重复注册根路由"
	StrPathExist     = "重复注册路由"
)

// RouterInterface 路由器的接口定义
type RouterInterface interface {
	// addRoute 添加路由
	addRoute(method string, path string, f4h HTTPHandleFunc, s5f4mw ...HTTPMiddleware)
	// findRoute 查找路由
	findRoute(method string, path string) *routeInfo
	// 服务器动前，将中间件扫描好，缓存到路由结点
	middlewareCache()
}

// router 路由器
type router struct {
	// m3routingTree 路由树，路由按请求方式分成多个路由树
	// map：Get => Get 的路由树；Post => Post 的路由树。
	m3routingTree map[string]*routingNode
}

func (p7this *router) addRoute(method string, path string, f4h HTTPHandleFunc, s5f4mw ...HTTPMiddleware) {
	if "" == path {
		panic(StrPathCannotBeEmpty)
	}
	if '/' != path[0] {
		panic(StrPathMustStartWithSlash)
	}
	if '/' == path[len(path)-1] && "/" != path {
		panic(StrPathCannotEndWithSlash)
	}

	// 按 http method 区分路由树
	p7node, ok := p7this.m3routingTree[method]
	if !ok {
		// 创建路由树根结点
		p7node = &routingNode{
			nodeType: nodeTypeStatic,
			part:     "/",
			path:     "/",
		}
		p7this.m3routingTree[method] = p7node
	}
	// 如果添加的根路由，处理完直接返回
	if "/" == path {
		if nil != p7node.f4handler {
			panic(StrRootNodeExist)
		}
		p7node.f4handler = f4h
		return
	}
	// 分段处理路由
	s5path := strings.Split(path[1:], "/")
	t4path := ""
	for _, part := range s5path {
		if "" == part {
			panic(StrPartCannotBeEmpty)
		}
		t4path += "/" + part
		t4p7child := p7node.findChild(part)
		if nil == t4p7child {
			t4p7child = p7node.createChild(part, t4path)
		} else {
			t4p7child.checkChild(part)
		}
		p7node = t4p7child
	}
	// 给路由添加处理方法
	if nil != p7node.f4handler {
		panic(StrPathExist)
	}
	p7node.f4handler = f4h
	// 给路由添加中间件
	p7node.s5f4middleware = s5f4mw
}

func (p7this *router) findRoute(method string, path string) *routeInfo {
	p7node, ok := p7this.m3routingTree[method]
	if !ok {
		return nil
	}

	if "/" == path {
		return &routeInfo{
			p7node: p7node,
		}
	}

	p7ri := &routeInfo{}
	s5path := strings.Split(path[1:], "/")
	for _, part := range s5path {
		p7node = p7node.matchChild(part)
		if nil == p7node {
			return nil
		}
		if p7node.nodeType == nodeTypeRegexp {
			s5res := p7node.p7regexp.FindStringSubmatch(part)
			if nil != s5res {
				p7ri.addPathParam(p7node.paramName, s5res[0])
			} else {
				p7ri.addPathParam(p7node.paramName, "")
			}
		} else if p7node.nodeType == nodeTypeParam {
			p7ri.addPathParam(p7node.paramName, part)
		} else if p7node.nodeType == nodeTypeAny {
			if nil == p7node.m3routingTree && nil == p7node.p7regexpChild && nil == p7node.p7paramChild && nil == p7node.p7anyChild {
				break
			}
		}
	}
	p7ri.p7node = p7node
	return p7ri
}

func (p7this *router) middlewareCache() {
	// 遍历每一个路由树
	for _, p7node := range p7this.m3routingTree {
		p7node.makeMiddlewareCache(nil)
	}
}

func newRouter() router {
	return router{
		m3routingTree: map[string]*routingNode{},
	}
}

// routeInfo 路由查询结果
type routeInfo struct {
	// 命中的路由结点
	p7node *routingNode
	// 提取出来的路径参数
	m3pathParam map[string]string
}

// addPathParam 添加路径参数
func (p7this *routeInfo) addPathParam(k string, v string) {
	if nil == p7this.m3pathParam {
		p7this.m3pathParam = map[string]string{k: v}
	}
	p7this.m3pathParam[k] = v
}
