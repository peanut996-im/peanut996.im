package http

import "github.com/gin-gonic/gin"

//Server ...
type Server struct {
	session *gin.Engine
	router  map[string]map[string]interface{}
}

//NewServer ...
func NewServer() *Server {
	return &Server{
		session: gin.Default(),
		router: map[string]map[string]interface{}{
			"GET":  nil,
			"POST": nil,
		},
	}
}

//Serve ...
func (s *Server) Serve() {
}
