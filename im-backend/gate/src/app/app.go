package app

import (
	"framework/cfgargs"
	"framework/logger"
	"gate/server"
	"github.com/gin-gonic/gin"
	"sync"
)

var (
	once sync.Once
	app  *App
)

type App struct {
	srv *server.Server
	cfg *cfgargs.SrvConfig
}

func GetApp() *App {
	once.Do(func() {
		a := &App{}
		app = a
	})
	return app
}

func (a *App) Init(cfg *cfgargs.SrvConfig) {
	gin.DefaultWriter = logger.MultiWriter(logger.DefLogger().GetLogWriters()...)
	//socket.io
	a.srv = server.NewServer()
	a.srv.Init(cfg)
}
