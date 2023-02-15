package pkg

import "strings"

const (
	//结点类型，越往下权重越大
	c1nodeTypeRoot   int = iota // 根结点
	c1nodeTypeAny               // '*' 匹配
	c1nodeTypeParam             // 路径参数
	c1nodeTypeStatic            //静态匹配
)

// 魔术变量
const c1anyStr = "*"

// node，树结点，私有
type node struct {

	// 结点类型
	nodeType int
	// 路由的一段
	pattren string
	// 子结点
	arr1p1children []*node

	// 如果这个结点是叶子结点，那么就可以调用路由对应的处理方法了
	hhFunc HTTPHandlerFunc
	// 判断是否匹配，如果是路径参数结点还是会顺带提取参数
	rmFunc func(path string, p1c *HTTPContext) bool
}

func newNode(path string) *node {
	if c1anyStr == path {
		return newAnyNode(path)
	} else if strings.HasPrefix(path, ":") {
		return newParamNode(path)
	}
	return newStaticNode(path)
}

func newRootNode(method string) *node {
	return &node{
		nodeType:       c1nodeTypeRoot,
		pattren:        method,
		arr1p1children: make([]*node, 0, 2),
		rmFunc: func(t1path string, p1c *HTTPContext) bool {
			// 根节点不允许匹配
			return false
		},
	}
}

func newStaticNode(path string) *node {
	return &node{
		nodeType:       c1nodeTypeStatic,
		pattren:        path,
		arr1p1children: make([]*node, 0, 2),
		rmFunc: func(t1path string, p1c *HTTPContext) bool {
			//静态匹配 不包括 "*"
			// 这里的Path 变量用到闭包的记忆效果
			return path == t1path && c1anyStr != t1path
		},
	}
}

func newAnyNode(path string) *node {
	return &node{
		nodeType: c1nodeTypeAny,
		pattren:  path,
		rmFunc: func(path string, p1c *HTTPContext) bool {
			// * 匹配，总是成功
			return true
		},
	}
}

func newParamNode(path string) *node {
	paramName := path[1:]
	return &node{
		nodeType: c1nodeTypeParam,
		pattren:  path,
		rmFunc: func(t1path string, p1c *HTTPContext) bool {
			// 路径参数，不包括 '*'
			isMatch := c1anyStr != t1path
			//顺便提取参数
			if isMatch && nil != p1c {
				p1c.PathParams[paramName] = t1path
			}
			return isMatch
		},
	}
}
