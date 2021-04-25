package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// type Route map[string]func(c *gin.Context)
// type Route map[func(c *gin.Context)] HTTPMethod
type Route struct {
	method  string
	path    string
	handler gin.HandlerFunc
}

// RouteNode distinguish routes with develop api.
type NodeRoute struct {
	path   string
	routes []*Route
}

// VerRoute distinguish routes with version.
// type VerRoute struct {
// 	Version   string
// 	DevRoutes []DevRoute
// }

//Server ...
type Server struct {
	session *gin.Engine
	routers []*NodeRoute
}

// NewServer ...
func NewServer() *Server {
	return &Server{
		session: gin.Default(),
		routers: []*NodeRoute{},
	}
}

func NewRoute(method, path string, handler gin.HandlerFunc) *Route {
	return &Route{
		method:  method,
		path:    path,
		handler: handler,
	}
}

func NewNodeRoute(path string, routers ...*Route) *NodeRoute {
	return &NodeRoute{
		path:   path,
		routes: routers,
	}
}
func (s *Server) AddNodeRoute(nodes ...*NodeRoute) {
	s.routers = append(s.routers, nodes...)
}

//Serve ...
func (s *Server) Serve(addr string) {
	s.mountRoutes()
	s.session.Use(cors.Default())
	err := s.session.Run(addr)
	if err != nil {
		panic(err)
	}
}

func (s *Server) mountRoutes() {
	router := s.session.Group("/")
	for _, node := range s.routers {
		group := router.Group(node.path)
		for _, route := range node.routes {
			methodMapper(group, route.method)(route.path, route.handler)
		}
	}
}

func methodMapper(group *gin.RouterGroup, method string) func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	switch method {
	case HTTP_METHOD_GET:
		return group.GET
	case HTTP_METHOD_POST:
		return group.POST
	case HTTP_METHOD_PUT:
		return group.PUT
	case HTTP_METHOD_DELETE:
		return group.DELETE
	case HTTP_METHOD_PATCH:
		return group.PATCH
	default:
		return group.Any
	}
}
