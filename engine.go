package aquarius

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

// 所有的group的engine都是最初创建的engine
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance
}

// 作为整个web服务的引擎，负责整个web服务的路由等逻辑
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // store all groups
}

// New is the constructor of Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		parent:      group,
		prefix:      group.prefix + prefix, // TODO 对于前缀的格式需要进行校验
		engine:      group.engine,
		middlewares: group.middlewares,
	}
	group.engine.groups = append(group.engine.groups, newGroup)
	return newGroup
}

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method, pattern string, handler HandlerFunc) {
	path := group.prefix + pattern
	log.Printf("Route %4s - %s", method, path)
	group.engine.router.addRoute(method, path, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (e *Engine) Run(addr string) (err error) {
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	return http.ListenAndServe(addr, e)
}

/*
TODO 确认这里的middleware能否按照正确逻辑进行组装
比如先添加 /hello/world/的中间件
再添加 /hello的中间件
遍历时会先讲 /hello/world的中间件添加到切片中
但是按照实际逻辑应该先添加/hello中的逻辑
不过按照使用逻辑，不应该添加这么添加

gin 的逻辑是，如果分别添加/hello/world和/hello两个group，那么这两个的middleware是相互独立的
*/
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	path := req.URL.Path
	middlewares := make([]HandlerFunc, 0)
	for _, g := range e.groups {
		if strings.HasPrefix(path, g.prefix) {
			middlewares = append(middlewares, g.middlewares...)
		}
	}
	c.handlers = middlewares
	e.router.handle(c)
}
