package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

var _ HTTPHandler = &HTTPHandlerTree{}

var arr1supportedMethod = [2]string{
	http.MethodGet,
	http.MethodPost,
}

// HTTPHandlerTree 基于前缀树实现路由处理
type HTTPHandlerTree struct {
	// 不同 method的根结点
	mapp1root map[string]*node
}

func NewHTTPHandlerTree() HTTPHandler {
	t1mapp1root := make(map[string]*node, len(arr1supportedMethod))
	for _, method := range arr1supportedMethod {
		t1mapp1root[method] = newRootNode(method)
	}
	return &HTTPHandlerTree{
		mapp1root: t1mapp1root,
	}
}

// HandlerHTTP HTTPHandler.HandlerHTTP
func (p1h *HTTPHandlerTree) HandlerHTTP(p1c *HTTPContext) {
	p1req := p1c.P1req
	fmt.Printf("HTTPHandlerTree, HandlerHTTP, p1req.Method: %s, p1req.URL.Path: %s\n", p1req.Method, p1req.URL.Path)

	hhFunc, err := p1h.findRoute(p1req.Method, p1req.URL.Path, p1c)
	if nil != err {
		p1c.P1resW.WriteHeader(http.StatusNotFound)
		_, _ = p1c.P1resW.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}

	hhFunc(p1c)
}

// findRoute 查询路由
func (p1h *HTTPHandlerTree) findRoute(method, path string, p1c *HTTPContext) (HTTPHandlerFunc, error) {
	t1path := strings.Trim(path, "/")
	arr1path := strings.Split(t1path, "/")

	p1nowNode, ok := p1h.mapp1root[method]
	if !ok {
		return nil, errors.New("method not")
	}

	for _, valpath := range arr1path {
		p1child := p1h.findMatchChild(p1nowNode, valpath, p1c)
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

	p1nowNode, ok := p1h.mapp1root[method]
	if !ok {
		return errors.New("method not supported")
	}

	for index, path := range arr1path {
		p1child := p1h.findMatchChild(p1nowNode, path, nil)
		if nil != p1child && c1nodeTypeAny != p1child.nodeType {
			//如果能找到子节点。移动指针到该节点，继续处理路由的下一段
			//这里额外的判断条件。是为了防止 '/user/*' 在 "/user/:id" 之前注册出问题
			p1nowNode = p1child
			continue
		}
		p1h.newSubTree(p1nowNode, arr1path[index:], hhFunc)
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
func (p1h *HTTPHandlerTree) findMatchChild(p1root *node, path string, p1c *HTTPContext) *node {
	arr1p1node := make([]*node, 0, 2)
	for _, p1child := range p1root.arr1p1children {
		if p1child.rmFunc(path, p1c) {
			arr1p1node = append(arr1p1node, p1child)
		}
	}

	if 0 == len(arr1p1node) {
		return nil
	}

	//根据结点权重排序
	sort.Slice(arr1p1node, func(i, j int) bool {
		return arr1p1node[i].nodeType < arr1p1node[j].nodeType
	})
	return arr1p1node[len(arr1p1node)-1]
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
