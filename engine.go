package aquarius

import (
	"log"
	"net/http"
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

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}
