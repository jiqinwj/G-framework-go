package pkg

import (
	"regexp"
	"strings"
)

const (
	// 静态路由
	nodeTypeStatic = iota
	// 通配符路由
	nodeTypeAny
	// 路径参数路由
	nodeTypeParam
	// 正则表达式路由
	nodeTypeRegexp
)

var (
	StrStaticChildExist = "重复注册静态路由"
	StrParamChildExist  = "重复注册路径参数路由"
	StrRegexpChildExist = "重复注册正则表达式路由"
	StrAnyChildExist    = "重复注册通配符路由"

	StrParamChildClashWithAnyChild    = "路径参数路由和通配符路由冲突"
	StrParamChildClashWithRegexpChild = "路径参数路由和正则表达式路由冲突"
	StrRegexpChildClashWithAnyChild   = "正则表达式路由和通配符路由冲突"
	StrRegexpChildClashWithParamChild = "正则表达式路由和路径参数路由冲突"
	StrAnyChildClashWithParamChild    = "通配符路由和路径参数路由冲突"
	StrAnyChildClashWithRegexpChild   = "通配符路由和正则表达式路由冲突"
)

// routingNode 路由结点
type routingNode struct {
	// nodeType 结点类型
	nodeType int
	// part 这个路由结点代表的那段路径
	part string
	// path 从根路由到这个路由结点的全路径
	path string

	// f4handler 命中路由之后的处理逻辑
	f4handler HTTPHandleFunc

	// m3routingTree 路由子树，子结点的 path => 子树根结点
	m3routingTree map[string]*routingNode
	// p7paramChild 路径参数结点
	p7paramChild *routingNode
	// paramName 路径参数路由和正则表达式路由，都会提取路由参数的名字
	paramName string
	// p7regexpChild 正则表达式结点
	p7regexpChild *routingNode
	// p7regexp 正则表达式
	p7regexp *regexp.Regexp
	// p7anyChild 通配符结点
	p7anyChild *routingNode

	// s5f4middleware 结点上注册的中间件
	s5f4middleware []HTTPMiddleware
	// s5f4middlewareCache 服务启动后，命中结点时，需要用到的所有中间件
	s5f4middlewareCache []HTTPMiddleware
}

// findChild 构建路由树时，查询子结点
func (p7this *routingNode) findChild(part string) *routingNode {
	// 找静态路由
	if nil != p7this.m3routingTree {
		t4p7node, ok := p7this.m3routingTree[part]
		if ok {
			return t4p7node
		}
	}
	// 找路径参数路由和正则表达式路由
	if ':' == part[0] {
		// 正则表达式用括号包裹
		t4regIndex1 := strings.Index(part, "(")
		t4regIndex2 := strings.Index(part, ")")
		if -1 != t4regIndex1 && -1 != t4regIndex2 && t4regIndex1 < t4regIndex2 {
			// 正则表达式路由
			return p7this.p7regexpChild
		} else {
			// 路径参数路由
			return p7this.p7paramChild
		}
	}
	// 找通配符路由
	if "*" == part {
		return p7this.p7anyChild
	}
	return nil
}

// checkChild 构建路由树时，校验子结点是否可以继续操作
func (p7this *routingNode) checkChild(part string) {
	// 这里需要校验路径参数路由和正则表达式路由是否冲突
	if ':' == part[0] {
		if p7this.part != part {
			panic(StrParamChildExist)
		}
	}
}

// createChild 构建路由树时，构造新的子结点
func (p7this *routingNode) createChild(part string, path string) *routingNode {
	if ':' == part[0] {
		t4regIndex1 := strings.Index(part, "(")
		t4regIndex2 := strings.Index(part, ")")
		if -1 != t4regIndex1 && -1 != t4regIndex2 && t4regIndex1 < t4regIndex2 {
			if nil != p7this.p7anyChild {
				panic(StrRegexpChildClashWithAnyChild)
			}
			if nil != p7this.p7paramChild {
				panic(StrRegexpChildClashWithParamChild)
			}
			if nil != p7this.p7regexpChild {
				panic(StrRegexpChildExist)
			}

			p7this.p7regexpChild = &routingNode{
				nodeType:  nodeTypeRegexp,
				part:      part,
				path:      path,
				paramName: part[1:t4regIndex1],
				p7regexp:  regexp.MustCompile(part[t4regIndex1+1 : t4regIndex2]),
			}
			return p7this.p7regexpChild
		} else {
			if nil != p7this.p7anyChild {
				panic(StrParamChildClashWithAnyChild)
			}
			if nil != p7this.p7paramChild {
				panic(StrParamChildExist)
			}
			if nil != p7this.p7regexpChild {
				panic(StrParamChildClashWithRegexpChild)
			}

			p7this.p7paramChild = &routingNode{
				nodeType:  nodeTypeParam,
				part:      part,
				path:      path,
				paramName: part[1:],
			}
			return p7this.p7paramChild
		}
	}
	if "*" == part {
		if nil != p7this.p7anyChild {
			panic(StrAnyChildExist)
		}
		if nil != p7this.p7paramChild {
			panic(StrAnyChildClashWithParamChild)
		}
		if nil != p7this.p7regexpChild {
			panic(StrAnyChildClashWithRegexpChild)
		}

		p7this.p7anyChild = &routingNode{
			nodeType: nodeTypeAny,
			part:     part,
			path:     path,
		}
		return p7this.p7anyChild
	}

	if nil == p7this.m3routingTree {
		p7this.m3routingTree = make(map[string]*routingNode)
	}
	_, ok := p7this.m3routingTree[part]
	if ok {
		panic(StrStaticChildExist)
	}

	p7this.m3routingTree[part] = &routingNode{
		nodeType: nodeTypeStatic,
		part:     part,
		path:     path,
	}
	return p7this.m3routingTree[part]
}

// makeMiddlewareCache 服务启动前，查询并缓存结点需要用到的所有中间件
func (p7this *routingNode) makeMiddlewareCache(s5f4mw []HTTPMiddleware) {
	t4s5f4mw := make([]HTTPMiddleware, 0, len(s5f4mw))
	// 上一层结点的中间件
	if nil != s5f4mw {
		t4s5f4mw = append(t4s5f4mw, s5f4mw...)
	}
	// 如果有通配符结点，则其他子结点需要把通配符结点上的中间件也加上
	if nil != p7this.p7anyChild {
		p7this.p7anyChild.makeMiddlewareCache(t4s5f4mw)
		if nil != p7this.p7anyChild.s5f4middleware {
			t4s5f4mw = append(t4s5f4mw, p7this.p7anyChild.s5f4middleware...)
		}
	}
	// 添加这个结点上的中间件
	if nil != p7this.s5f4middleware {
		t4s5f4mw = append(t4s5f4mw, p7this.s5f4middleware...)
	}

	// 如果这个结点有处理方法，那么这个结点就不是中间结点而是有效的路由结点，需要缓存中间件结果
	if nil != p7this.f4handler {
		p7this.s5f4middlewareCache = make([]HTTPMiddleware, 0, len(t4s5f4mw))
		p7this.s5f4middlewareCache = append(p7this.s5f4middlewareCache, t4s5f4mw...)
	}
	// 处理其余类型的子结点
	if nil != p7this.p7regexpChild {
		p7this.p7regexpChild.makeMiddlewareCache(t4s5f4mw)
	}
	if nil != p7this.p7paramChild {
		p7this.p7paramChild.makeMiddlewareCache(t4s5f4mw)
	}
	for _, p7childNode := range p7this.m3routingTree {
		p7childNode.makeMiddlewareCache(t4s5f4mw)
	}
}

// matchChild 查询路由时，匹配子结点
func (p7this *routingNode) matchChild(part string) *routingNode {
	// 这里的查询优先级可以根据需要进行调整
	// 先查询静态路由
	if nil != p7this.m3routingTree {
		p7node, ok := p7this.m3routingTree[part]
		if ok {
			return p7node
		}
	}
	// 然后依次查询，正则表达式路由、路径参数路由、通配符路由
	if nil != p7this.p7regexpChild {
		return p7this.p7regexpChild
	} else if nil != p7this.p7paramChild {
		return p7this.p7paramChild
	} else if nil != p7this.p7anyChild {
		return p7this.p7anyChild
	}
	return nil
}
