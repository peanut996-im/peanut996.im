package http

import (
	"fmt"
	"framework/api"
	"framework/cfgargs"
	"github.com/gin-gonic/gin"
)

type Route struct {
	method  string
	path    string
	handler gin.HandlerFunc
}

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
	cfg     *cfgargs.SrvConfig
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

//Init Read HTTP configuration: cross-domain, signature, Release
func (s *Server) Init(cfg *cfgargs.SrvConfig) {
	if cfg.HTTP.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	if cfg.HTTP.Cors {
		s.session.Use(CORS())
	}
	if cfg.HTTP.Sign {
		s.session.Use(CheckSign(cfg))
	}
	s.cfg = cfg
}

//Run Route and start HTTP based on the port of the yaml file.
func (s *Server) Run() error {
	s.mountRoutes()
	err := s.session.Run(fmt.Sprintf(":%v", s.cfg.HTTP.Port))
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Use(middlewares ...gin.HandlerFunc) {
	s.session.Use(middlewares...)
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
	case api.HTTPMethodGet:
		return group.GET
	case api.HTTPMethodPost:
		return group.POST
	case api.HTTPMethodPut:
		return group.PUT
	case api.HTTPMethodDelete:
		return group.DELETE
	case api.HTTPMethodPatch:
		return group.PATCH
	default:
		return group.Any
	}
}
