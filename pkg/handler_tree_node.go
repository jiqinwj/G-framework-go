package pkg

// node，树结点，私有
type node struct {
	// 路由的一段
	path string
	// 子结点
	arr1p1children []*node

	// 如果这个结点是叶子结点，那么就可以调用路由对应的处理方法了
	hhFunc HTTPHandlerFunc
}

func newNode(path string) *node {
	return &node{
		path:           path,
		arr1p1children: make([]*node, 0, 2),
	}
}
