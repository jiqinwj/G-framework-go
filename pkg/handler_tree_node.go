package pkg

// node 树结点，私有方法 小写 go的特性
type node struct {
	//路由的路径
	path string
	//子结点
	arr1n2children []*node

	//如果这个结点是叶子结点。那么就可以调用对应路由的处理方法
	hhFunc HTTPHandlerFunc
}

func newNode(path string) *node {
	return &node{
		path:           path,
		arr1n2children: make([]*node, 0, 2),
	}
}
