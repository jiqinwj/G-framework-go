package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// 这种方式 是为了避免 某个实现结构体 没有实现某个接口的方法
var _ HTTPHandler = &HTTPHandlerTree{}

// HTTPHandlerTree 基于一个前缀树实现路由处理
type HTTPHandlerTree struct {
	// 根节点
	n2root *node
}

func NewHTTPHandlerTree() HTTPHandler {
	return &HTTPHandlerTree{
		n2root: &node{},
	}
}

// HandlerHTTP 其实这个方法 是实现 这个HTTPHandler 接口里面的 HandlerHTTP 方法
func (n2h *HTTPHandlerTree) HandlerHTTP(c *HTTPContext) {
	n2req := c.N1req
	//打印下路由请求方法，路径
	fmt.Printf("HTTPHandlerTree,HandlerHTTP,n2req.Method: %s,n2req.URL.Path:%s\n", n2req.Method, n2req.URL.Path)
	//寻找路由了，找路了，很迷茫
	hhFunc, err := n2h.findRoute(n2req.URL.Path)
	if nil != err {
		c.N1resW.WriteHeader(http.StatusNotFound)
		_, _ = c.N1resW.Write([]byte(fmt.Sprintf("%s", err)))
	}
	//执行路由方法
	hhFunc(c)
}

// findRoute 查找路由了
func (n2h *HTTPHandlerTree) findRoute(path string) (HTTPHandlerFunc, error) {
	t1path := strings.Trim(path, "/")
	arr1path := strings.Split(t1path, "/")
	n2nowNode := n2h.n2root
	for _, valpath := range arr1path {
		n2child := n2h.findMatchChildV2(n2nowNode, valpath)
		if nil == n2child {
			return nil, errors.New("route not found")
		}
		n2nowNode = n2child
	}
	//防止访问到非叶子结点上，比如 注册了 '/user/info' 但是访问了 '/user'
	if nil == n2nowNode.hhFunc {
		return nil, errors.New("route not found")
	}
	return n2nowNode.hhFunc, nil
}

// findMatchChild 查找当前结点的子节点是否存在路由的分段
func (n2h *HTTPHandlerTree) findMatchChild(n2root *node, path string) *node {
	for _, n2child := range n2root.arr1n2children {
		if path == n2child.path {
			return n2child
		}
	}
	return nil
}

// findMatchChildV2 查找当前结点的子节点是否存在路由的另外一条路
func (n2h *HTTPHandlerTree) findMatchChildV2(n2Root *node, path string) *node {
	var t1p1node *node
	for _, n2child := range n2Root.arr1n2children {
		//查找到了直接返回得了
		if path == n2child.path && "*" != n2child.path {
			return n2child
		}
		//没有就继续查找
		if "*" == n2child.path {
			t1p1node = n2child
		}
	}
	return t1p1node
}

// RegisteRoute HTTPHandler.HTTPRoute.RegisteRoute
func (n2h *HTTPHandlerTree) RegisteRoute(method string, pattern string, hhFunc HTTPHandlerFunc) error {
	//检查格式 自定义匹配模式
	err := n2h.checkPattern(pattern)
	if err != nil {
		return err
	}

	// 清理路由前后的"/",然后把路由分段
	t1pattern := strings.Trim(pattern, "/")
	arr1path := strings.Split(t1pattern, "/")
	//指针指向当前操作的结点
	n2nowNode := n2h.n2root
	//依次处理路由的每一段
	for index, path := range arr1path {
		n2child := n2h.findMatchChild(n2nowNode, path)
		if nil != n2child {
			//如果能找到子节点，移动指针到该节点，继续处理路由的下一段
			n2nowNode = n2child
			continue
		}
		// 如果找不到子节点。说明需要创建新的子树
		n2h.newSubTree(n2nowNode, arr1path[index:], hhFunc)
		break
	}

	//如果能找到子节点但是没有创建新的子树
	//证明是短路由 比方先注册 /user/info .再注册/user
	n2nowNode.hhFunc = hhFunc

	return nil

}

// checkPattern 校验路由
// 对于 "*" 只允许 '/*' 或者 "/user/*" 格式存在
// 不允许 "/user*" 或者 "/user/*/info" 存在
func (n2h *HTTPHandlerTree) checkPattern(pattern string) error {
	index := strings.Index(pattern, "*")
	if index == 0 {
		return errors.New("route pattern is errro,index==0")
	} else if index > 0 {
		//判断 "*" 是不是最后一个字符
		if len(pattern)-1 != index {
			return errors.New("route pattern is error,len(pattern) - 1!=index")
		}
		//判断不是 ".../*" 格式
		if '/' != pattern[index-1] {
			return errors.New("route pattern is error,'/' !=patter[index-1] ")
		}
	}
	return nil
}

// newSubTree 创建子树
func (n2h *HTTPHandlerTree) newSubTree(n2root *node, arr1path []string, hhFunc HTTPHandlerFunc) {

	n2nowNode := n2root
	for _, path := range arr1path {
		//创建新的节点
		n2newNode := newNode(path)
		n2nowNode.arr1n2children = append(n2nowNode.arr1n2children, n2newNode)
		//移动指针到新节点，继续处理路由的下一段
		n2nowNode = n2newNode
	}
	// 子树构造完成，在叶子节点上绑定路由对应的处理方法
	n2nowNode.hhFunc = hhFunc

}
