package pkg

import (
	"regexp"
	"strings"
)

const (
	// 静态路由
	nodeTypeStatic = iota
	//通配符路由
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

type routingNode struct {
	// 结点类型
	nodeType int
	// part 这个路由结点代表的那段路径
	part string
	// path 从根路由到这个路由结点的全路径
	path string

	// 命中路由之后的处理逻辑
	f4handler HTTPHandlerFunc

	//路由子树，子结点的 path =>子树根结点
	m3routingTree map[string]*routingNode
	// 路径参数结点
	p7paramChild *routingNode
	// 路径参数由和正则表达式路由，都会提取路由参数的名字
	paramName string
	// 正则表达式结点
	p7regexpChild *routingNode
	// 正则表达式
	p7regexp *regexp.Regexp
	// 通配符结点
	p7anyChild *routingNode

	// 结点上注册的中间件
	s5f4middleware []HTTPMiddleware
	// 服务启动后，命中结点时，需要用到的所有中间件
	s5f4middlewareCache []HTTPMiddleware
}

// 构建路由树时，查询子节点
func (p7this *routingNode) findChild(part string) *routingNode {
	// 找静态路由
	if nil != p7this.m3routingTree {
		t4p7node, ok := p7this.m3routingTree[part]
		if ok {
			return t4p7node
		}
	}

	//找路径参数路由和正则表达式路由
	if ':' == part[0] {
		//正则表达式用括号包裹
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

// 构建路由树时，构建新的子结点
func (p7this *routingNode) createChild(part string, path string) *routingNode {
	// 注册正则匹配路由
	if ':' == part[0] {
		t4regIndex1 := strings.Index(part, "(")
		t4regIndex2 := strings.Index(part, ")")
		if -1 != t4regIndex1 && -1 != t4regIndex2 && t4regIndex1 < t4regIndex2 {
			// 各种规则冲突判断
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
		}
	}

	//注册通配符路由
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

	//如果开始都是根节点 直接构建一个
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

// 构建路由树时，校验子结点是否可以继续操作
func (p7this *routingNode) checkChild(part string) {
	// 这里需要校验路径参数路由和正则表达式是否冲突
	if ':' == part[0] {
		if p7this.part != part {
			panic(StrParamChildExist)
		}
	}
}
