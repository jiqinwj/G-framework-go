package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var _ HTTPHandler = &HTTPHandlerTree{}

// HTTPHandlerTree 基于前缀树实现路由处理
type HTTPHandlerTree struct {
	// 根结点
	p1root *node
}

func NewHTTPHandlerTree() HTTPHandler {
	return &HTTPHandlerTree{
		p1root: &node{},
	}
}

// HandlerHTTP HTTPHandler.HandlerHTTP
func (p1h *HTTPHandlerTree) HandlerHTTP(c *HTTPContext) {
	p1req := c.P1req
	fmt.Printf("HTTPHandlerTree, HandlerHTTP, p1req.Method: %s, p1req.URL.Path: %s\n", p1req.Method, p1req.URL.Path)

	hhFunc, err := p1h.findRoute(p1req.URL.Path)
	if nil != err {
		c.P1resW.WriteHeader(http.StatusNotFound)
		_, _ = c.P1resW.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}

	hhFunc(c)
}

// findRoute 查询路由
func (p1h *HTTPHandlerTree) findRoute(path string) (HTTPHandlerFunc, error) {
	t1path := strings.Trim(path, "/")
	arr1path := strings.Split(t1path, "/")
	p1nowNode := p1h.p1root
	for _, valpath := range arr1path {
		p1child := p1h.findMatchChildV2(p1nowNode, valpath)
		if nil == p1child {
			return nil, errors.New("route not found")
		}
		p1nowNode = p1child
	}
	// 防止访问到非叶子结点上。比如，注册了 `/user/info` 但是访问 `/user`。
	if nil == p1nowNode.hhFunc {
		return nil, errors.New("route not found")
	}
	return p1nowNode.hhFunc, nil
}

// findMatchChildV2 查找当前结点的子结点是否存在路由的分段
func (p1h *HTTPHandlerTree) findMatchChildV2(p1root *node, path string) *node {
	var t1p1node *node
	for _, p1child := range p1root.arr1p1children {
		if path == p1child.path && "*" != p1child.path {
			return p1child
		}
		if "*" == p1child.path {
			t1p1node = p1child
		}
	}
	return t1p1node
}

// RegisteRoute HTTPHandler.HTTPRoute.RegisteRoute
func (p1h *HTTPHandlerTree) RegisteRoute(method string, pattern string, hhFunc HTTPHandlerFunc) error {
	err := p1h.checkPattern(pattern)
	if nil != err {
		return err
	}

	// 清理路由前后的 `/`，然后把路由分段
	t1pattern := strings.Trim(pattern, "/")
	arr1path := strings.Split(t1pattern, "/")
	fmt.Printf("路由数据:%v \n", arr1path)

	// 指针指向当前操作的结点
	p1nowNode := p1h.p1root
	// 依次处理路由的每一段
	for index, path := range arr1path {
		p1child := p1h.findMatchChild(p1nowNode, path)
		if nil != p1child {
			// 如果能找到子结点，移动指针到该子结点，继续处理路由的下一段
			p1nowNode = p1child
			continue
		}
		// 如果找不到子结点，说明需要创建新的子树
		fmt.Printf("创建子树节点：%v \n", arr1path[index:])
		p1h.newSubTree(p1nowNode, arr1path[index:], hhFunc)
		break
	}

	// 如果能找到子结点但是没有创建新的子树。
	// 说明注册的是短路由，比如先注册 /user/info，再注册 /user。
	p1nowNode.hhFunc = hhFunc
	return nil
}

// checkPattern 校验路由
// 对于 `*`，只允许 `/*` 或 `/user/*`，格式存在。
// 不允许 `/user*` 或 `/user/*/info` 存在
func (p1h *HTTPHandlerTree) checkPattern(pattern string) error {
	index := strings.Index(pattern, "*")
	if index == 0 {
		return errors.New("route pattern is error, index == 0")
	} else if index > 0 {
		// 判断 `*` 是不是最后一个字符
		if len(pattern)-1 != index {
			return errors.New("route pattern is error, len(pattern) - 1 != index")
		}
		// 判断是不是 `.../*` 格式
		if '/' != pattern[index-1] {
			return errors.New("route pattern is error, '/' != pattern[index-1]")
		}
	}
	return nil
}

// findMatchChild 查找当前结点的子结点是否存在路由的分段
func (p1h *HTTPHandlerTree) findMatchChild(p1root *node, path string) *node {
	for _, p1child := range p1root.arr1p1children {
		if path == p1child.path {
			return p1child
		}
	}
	return nil
}

// newSubTree 创建子树
func (p1h *HTTPHandlerTree) newSubTree(p1root *node, arr1path []string, hhFunc HTTPHandlerFunc) {
	p1nowNode := p1root
	for _, path := range arr1path {
		// 创建新的结点
		p1newNode := newNode(path)
		p1nowNode.arr1p1children = append(p1nowNode.arr1p1children, p1newNode)
		// 移动指针到新结点，继续处理路由的下一段
		p1nowNode = p1newNode
	}
	// 子树构造完成，在叶子结点上绑定路由对应的处理方法
	p1nowNode.hhFunc = hhFunc
}
