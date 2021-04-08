package aquarius

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type router struct {
	tries    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		tries:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

/*
路由可能的情况
/
/hello/
world
/hello/:name/
/hello/world/*msg
*/
func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
	root, ok := r.tries[method]
	if !ok {
		root = new(node)
		r.tries[method] = root
	}
	parts := parsePattern(pattern)
	// 向trie树种添加路由
	root.insert(pattern, parts, 0)
	// 向handlerMap中添加对应模式串和方法的handler
	r.handlers[key] = handler
}

func (r *router) getRouter(method, pattern string) (*node, map[string]string) {
	root, ok := r.tries[method]
	if !ok {
		return nil, nil
	}
	searchParts := parsePattern(pattern)
	n := root.search(searchParts, 0)
	if n == nil {
		return nil, nil
	}
	params := make(map[string]string)
	parts := parsePattern(n.pattern)
	for i, part := range parts {
		if part[0] == ':' {
			params[part[1:]] = searchParts[i]
		} else if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(searchParts[i:], "/")
			break
		}
	}
	return n, params
}

func (r *router) handle(c *Context) {
	root, params := r.getRouter(c.Method, c.Path)
	c.Params = params
	if root != nil {
		key := c.Method + "-" + root.pattern
		if handler, ok := r.handlers[key]; ok {
			handler(c)
			return
		}
		log.Println(fmt.Sprintf("the handler should not be nil, because there is the mathched tire node, [method: %s, path: %s, trie.pattern: %s]", c.Method, c.Path, root.pattern))
	}
	c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}
