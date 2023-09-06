package web

import (
	"fmt"
	"regexp"
	"strings"
)

type PATHTYPE int

const (
	STARPATH PATHTYPE = iota
	PARAMSPATH
	GZPATH
	FULLPATH
)

type node struct {
	path string
	// 全匹配
	children map[string]*node
	// 通配符(*)匹配
	starChildren *node
	handle       HandleFunc
	// 参数匹配
	paramsChildren *node
	// 正则匹配
	gzChildren *node
	reg        *regexp.Regexp
	nodeType   PATHTYPE
	matchRoute string

	//middleware
	ms []Middleware
}

type nodeInfo struct {
	*node
	params map[string]string
	ms     []Middleware
}

type router struct {
	trees map[string]*node
}

func NewRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

func (n *node) child(seg string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramsChildren != nil {
			return n.paramsChildren, true, true
		}
		if n.gzChildren != nil {
			if n.gzChildren.reg.MatchString(seg) {
				return n.gzChildren, true, true
			}
			return nil, false, true
		}
		if n.starChildren != nil {
			return n.starChildren, false, n.starChildren != nil
		}
		if n.nodeType == STARPATH {
			return n, false, true
		}
	}

	nextNode, ok := n.children[seg]
	if !ok {
		if n.paramsChildren != nil {
			return n.paramsChildren, true, true
		}
		if n.gzChildren != nil {
			if n.gzChildren.reg.MatchString(seg) {
				return n.gzChildren, true, true
			}
		}
		if n.starChildren != nil {
			return n.starChildren, false, true
		}
		if n.nodeType == STARPATH {
			return n, false, true
		}
		return nil, false, false
	}
	return nextNode, false, true
}
func (h *router) findRoute(method string, path string) (*nodeInfo, bool) {
	_, ok := checkPath(path)
	if !ok {
		return nil, ok
	}
	node, ok := h.trees[method]
	if !ok {
		panic(fmt.Sprintf("web：未找到对应请求类型 %s", method))
	}
	if path == "/" {
		return &nodeInfo{
			node: node,
		}, true
	}
	path = path[1:]
	segs := strings.Split(path, "/")
	var params map[string]string
	var supMiddleware []Middleware
	for _, seg := range segs {
		if seg == "" {
			return nil, false
		}
		nextNode, isParams, ok := node.child(seg)
		if !ok {
			return nil, false
		}
		if isParams {
			if params == nil {
				params = map[string]string{}
			}
			params[nextNode.path] = seg
		}
		if len(node.ms) > 0 {
			supMiddleware = append(supMiddleware, node.ms...)
		}
		node = nextNode
	}
	return &nodeInfo{node, params, supMiddleware}, true
}

func checkPath(path string) (string, bool) {
	if path == "" {
		return "web：路由路径不能未空", false
	}
	// 必须以'/'开头
	if path[0] != '/' {
		return "web：路由路径必须以'/'开始", false
	}
	if path != "/" && path[len(path)-1] == '/' {
		return "web：路由路径不能以'/'结束", false
	}
	return "", true
}

func (n *node) findOrCreateNode(seg string) *node {
	if seg == "" {
		panic("web：路由路径不支持'//'")
	}
	// 参数匹配
	if seg[0] == ':' {
		if n.starChildren != nil {
			panic("web：已有通配符路由匹配")
		}
		if n.gzChildren != nil {
			panic("web：已有正则路由匹配")
		}
		// 正则匹配
		seg = seg[1:]
		segs := strings.SplitN(seg, "(", 2)
		if len(segs) == 2 {
			expr := segs[1]
			if strings.HasSuffix(expr, ")") {
				if n.gzChildren == nil {
					reg := regexp.MustCompile(expr[:len(expr)-1])
					gzNode := createNewNode(segs[0])
					gzNode.reg = reg
					gzNode.nodeType = GZPATH
					n.gzChildren = gzNode
				}
				return n.gzChildren
			}
		}
		if n.paramsChildren == nil {
			paramsNode := createNewNode(seg)
			paramsNode.nodeType = PARAMSPATH
			n.paramsChildren = paramsNode
		}
		return n.paramsChildren
	}

	// 通配符匹配
	if seg == "*" {
		if n.paramsChildren != nil {
			panic("web：已有通参数路由匹配")
		}
		if n.starChildren == nil {
			starNode := createNewNode(seg)
			starNode.nodeType = STARPATH
			n.starChildren = starNode
		}
		if n.gzChildren != nil {
			panic("web：已有正则路由匹配")
		}
		return n.starChildren
	}
	// 静态匹配
	if n.children == nil {
		n.children = map[string]*node{}
	}
	return n.findChild(seg)
}

func (h *router) AddRoute(method string, path string, handleFunc HandleFunc, ms ...Middleware) {
	msg, ok := checkPath(path)
	if !ok {
		panic(msg)
	}
	root, ok := h.trees[method]
	if !ok {
		h.trees[method] = createNewNode("/")
		root = h.trees[method]
	}
	if path == "/" {
		if root.handle != nil {
			panic("web：路由路径不能重复")
		}
		root.handle = handleFunc
		return

	}
	path = path[1:]
	segments := strings.Split(path, "/")
	for _, seg := range segments {
		nextNode := root.findOrCreateNode(seg)
		root = nextNode
	}
	if root.handle != nil {
		panic("web：路由路径不能重复")
	}
	root.handle = handleFunc
	root.matchRoute = path
	root.ms = ms
}

func (root *node) findChild(seg string) *node {
	childrenNode, ok := root.children[seg]
	if !ok {
		newNode := createNewNode(seg)
		newNode.nodeType = FULLPATH
		root.children[seg] = newNode
		return newNode
	}
	return childrenNode
}

func createNewNode(seg string) *node {
	return &node{
		path:     seg,
		nodeType: FULLPATH,
	}
}
