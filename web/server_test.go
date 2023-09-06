package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"regexp"
	"testing"
)

func genTestRoute(mockHandleFunc HandleFunc) *router {
	testCases := []struct {
		methodName string
		path       string
	}{
		{
			methodName: http.MethodGet,
			path:       "/user/login",
		},
		{
			methodName: http.MethodGet,
			path:       "/user",
		},
		{
			methodName: http.MethodPost,
			path:       "/user/register",
		},
		{
			methodName: http.MethodPost,
			path:       "/login/*",
		},
		{
			methodName: http.MethodPost,
			path:       "/nameId/:id/:name",
		},
		{
			methodName: http.MethodPut,
			path:       "/getUserInfo/:id(^[0-9]*$)/:userName/*",
		},
	}
	route := NewRouter()
	for _, tc := range testCases {
		route.AddRoute(tc.methodName, tc.path, mockHandleFunc)
	}
	return route
}

func TestRouter_AddRoute(t *testing.T) {
	var mockHandleFunc HandleFunc = func(ctx *Context) {}
	testRouter := genTestRoute(mockHandleFunc)
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path: "/",
				children: map[string]*node{
					"user": &node{
						path:     "user",
						nodeType: FULLPATH,
						handle:   mockHandleFunc,
						children: map[string]*node{
							"login": &node{
								path:     "login",
								handle:   mockHandleFunc,
								nodeType: FULLPATH,
							},
						},
					},
				},
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"user": &node{
						path:     "user",
						nodeType: FULLPATH,
						children: map[string]*node{
							"register": &node{
								path:     "register",
								nodeType: FULLPATH,
								handle:   mockHandleFunc,
							},
						},
					},
					"login": &node{
						path:     "login",
						nodeType: FULLPATH,
						starChildren: &node{
							path:     "*",
							nodeType: STARPATH,
							handle:   mockHandleFunc,
						},
					},
					"nameId": &node{
						path:     "nameId",
						nodeType: FULLPATH,
						paramsChildren: &node{
							path:     "id",
							nodeType: PARAMSPATH,
							paramsChildren: &node{
								path:     "name",
								handle:   mockHandleFunc,
								nodeType: PARAMSPATH,
							},
						},
					},
				},
			},
			http.MethodPut: &node{
				path: "/",
				children: map[string]*node{
					"getUserInfo": &node{
						path:     "getUserInfo",
						nodeType: FULLPATH,
						gzChildren: &node{
							path:     "id",
							nodeType: GZPATH,
							reg:      regexp.MustCompile("(^[0-9]*$)"),
							paramsChildren: &node{
								path:     "userName",
								nodeType: PARAMSPATH,
								starChildren: &node{
									path:     "*",
									nodeType: STARPATH,
									handle:   mockHandleFunc,
								},
							},
						},
					},
				},
			},
		},
	}
	msg, isequal := wantRouter.equal(testRouter)
	assert.True(t, isequal, msg)

	// 校验功能测试
	routePanic := NewRouter()
	assert.PanicsWithValue(t, "web：路由路径不能未空", func() {
		routePanic.AddRoute(http.MethodGet, "", mockHandleFunc)
	})
	assert.PanicsWithValue(t, "web：路由路径必须以'/'开始", func() {
		routePanic.AddRoute(http.MethodGet, "login", mockHandleFunc)
	})
	assert.PanicsWithValue(t, "web：路由路径不能以'/'结束", func() {
		routePanic.AddRoute(http.MethodGet, "/login/user/", mockHandleFunc)
	})
	assert.PanicsWithValue(t, "web：路由路径不支持'//'", func() {
		routePanic.AddRoute(http.MethodGet, "/login///user", mockHandleFunc)
	})
	assert.PanicsWithValue(t, "web：路由路径不能重复", func() {
		routePanic.AddRoute(http.MethodGet, "/login/user", mockHandleFunc)
		routePanic.AddRoute(http.MethodGet, "/login/user", mockHandleFunc)
	})
	assert.PanicsWithValue(t, "web：已有通配符路由匹配", func() {
		routePanic.AddRoute(http.MethodGet, "/login/*", mockHandleFunc)
		routePanic.AddRoute(http.MethodGet, "/login/:id", mockHandleFunc)
	})
	routePanic = NewRouter()
	assert.PanicsWithValue(t, "web：已有通参数路由匹配", func() {
		routePanic.AddRoute(http.MethodGet, "/login/:id", mockHandleFunc)
		routePanic.AddRoute(http.MethodGet, "/login/*", mockHandleFunc)
	})
}

func (r *router) equal(y *router) (string, bool) {
	if len(r.trees) != len(y.trees) {
		return fmt.Sprintf("web：路由方法数量不匹配"), false
	}
	for k, t := range r.trees {
		yt, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("web：路由方法数量不匹配,没有%s方法", k), false
		}
		msg, isequal := t.equal(yt)
		if !isequal {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(m *node) (string, bool) {
	if n.path != m.path {
		return fmt.Sprintf("web：路径不匹配 目标路径%s 原路径%s", n.path, m.path), false
	}
	if n.nodeType != m.nodeType {
		return fmt.Sprintf("web：节点类型不匹配 目标节点类型%d 原节点类型%d", n.nodeType, m.nodeType), false
	}
	nHandle := reflect.ValueOf(n.handle)
	mHanlde := reflect.ValueOf(m.handle)
	if mHanlde != nHandle {
		return fmt.Sprintf("web：handler 方法不匹配"), false
	}

	if len(n.children) != len(m.children) {
		return fmt.Sprintf("web：子节点数量不匹配"), false
	}
	for k, nc := range n.children {
		mc, ok := m.children[k]
		if !ok {
			return fmt.Sprintf("web：目标节点缺少路径%s", k), false
		}
		msg, ok := nc.equal(mc)
		if !ok {
			return msg, false
		}
	}
	if n.starChildren != nil {
		msg, ok := n.starChildren.equal(m.starChildren)
		if !ok {
			return msg, false
		}
	}
	if n.paramsChildren != nil {
		msg, ok := n.paramsChildren.equal(m.paramsChildren)
		if !ok {
			return msg, false
		}
	}
	return "", true
}

func TestRouter_findRoute(t *testing.T) {
	var mockHandleFunc HandleFunc = func(ctx *Context) {}
	testRouter := genTestRoute(mockHandleFunc)
	testCases := []struct {
		name       string
		methodName string
		path       string
		wantFound  bool
		wantParams map[string]string
		wantNode   *node
	}{
		{
			name:       "findGetNode",
			methodName: http.MethodGet,
			path:       "/",
			wantFound:  true,
			wantNode: &node{
				path: "/",
				children: map[string]*node{
					"user": &node{
						path:   "user",
						handle: mockHandleFunc,
						children: map[string]*node{
							"login": &node{
								path:   "login",
								handle: mockHandleFunc,
							},
						},
					},
				},
			},
		},
		{
			name:       "findPostNode",
			methodName: http.MethodPost,
			wantFound:  true,
			path:       "/user",
			wantNode: &node{
				path: "user",
				children: map[string]*node{
					"register": &node{
						path:   "register",
						handle: mockHandleFunc,
					},
				},
			},
		},
		{
			name:       "not found",
			methodName: http.MethodPost,
			wantFound:  false,
			path:       "/user/notFound",
		},
		{
			name:       "findStart",
			methodName: http.MethodPost,
			wantFound:  true,
			path:       "/login/123",
			wantNode: &node{
				path:   "*",
				handle: mockHandleFunc,
			},
		},
		{
			name:       "findParam",
			methodName: http.MethodPost,
			wantFound:  true,
			path:       "/nameId/123",
			wantNode: &node{
				path: "id",
				paramsChildren: &node{
					path:   "name",
					handle: mockHandleFunc,
				},
			},
			wantParams: map[string]string{
				"id": "123",
			},
		},
		{
			name:       "findParams",
			methodName: http.MethodPost,
			wantFound:  true,
			path:       "/nameId/123/xiejunze",
			wantNode: &node{
				path:   "name",
				handle: mockHandleFunc,
			},
			wantParams: map[string]string{
				"id":   "123",
				"name": "xiejunze",
			},
		},
		{
			name:       "findRz",
			methodName: http.MethodPut,
			wantFound:  true,
			path:       "/getUserInfo/123/AL/test",
			wantNode: &node{
				path:   "*",
				handle: mockHandleFunc,
			},
			wantParams: map[string]string{
				"id":       "123",
				"userName": "AL",
			},
		},
		{
			name:       "findStart",
			methodName: http.MethodPost,
			wantFound:  true,
			path:       "/login/123/456",
			wantNode: &node{
				path:   "*",
				handle: mockHandleFunc,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node, ok := testRouter.findRoute(tc.methodName, tc.path)
			assert.Equal(t, tc.wantFound, ok)
			if ok {
				msg, ok := node.equal(tc.wantNode)
				assert.True(t, ok, msg)
				assert.Equal(t, node.params, tc.wantParams)
			}
		})
	}
}
