package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//Server ...
type Server struct {
	session *gin.Engine
	router  map[string]map[string]http.Handler
}

// NewServer ...
func NewServer() *Server {
	return &Server{
		session: gin.Default(),
		router: map[string]map[string]http.Handler{
			"GET":  nil,
			"POST": nil,
		},
	}
}

//Serve ...
func (s *Server) Serve() {
}
