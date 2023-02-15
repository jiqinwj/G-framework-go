package pkg

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHTTPHandlerTree_RegisteRoute(t *testing.T) {
	// 创建接口实列，记得断言成测试的对象
	p1h := NewHTTPHandlerTree().(*HTTPHandlerTree)
	// 检查支持的 HTTP 方法数量是否正确
	assert.Equal(t, len(arr1supportedMethod), len(p1h.mapp1root))

	// 添加 PUT /user
	err := p1h.RegisteRoute(http.MethodPut, "/user", func(p1c *HTTPContext) {})
	// 检查是否抛出异常
	assert.Equal(t, errors.New("method not supported"), err)

	//添加Get /user
	err = p1h.RegisteRoute(http.MethodGet, "/user", func(p1c *HTTPContext) {})
	assert.Nil(t, err)

	//检查 Get 组/结点，子节点数量
	t1getNode := p1h.mapp1root[http.MethodGet]
	assert.Equal(t, 1, len(t1getNode.arr1p1children))
	// 检查GET组/user 结点
	t1userNode := t1getNode.arr1p1children[0]
	assert.NotNil(t, t1userNode)
	assert.Equal(t, "user", t1userNode.pattren)
	assert.Empty(t, t1userNode.arr1p1children)
	assert.NotNil(t, t1userNode.hhFunc)

	////添加 GET /user/info
	err = p1h.RegisteRoute(http.MethodGet, "/user/info", func(c *HTTPContext) {})
	assert.Nil(t, err)
	////检查 GET组 /user 结点，子节点数量
	assert.Equal(t, 1, len(t1userNode.arr1p1children))
	////检查 GET组 /user/info 结点
	t1userInfoNode := t1userNode.arr1p1children[0]
	assert.NotNil(t, t1userInfoNode)
	assert.Equal(t, "info", t1userInfoNode.pattren)
	assert.Empty(t, t1userInfoNode.arr1p1children)
	assert.NotNil(t, t1userInfoNode.hhFunc)

}
