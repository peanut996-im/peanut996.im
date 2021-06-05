package app

import (
	"framework/cfgargs"
	"framework/db"
	"framework/logger"
	"github.com/gin-gonic/gin"
	"logic/server"
	"sync"
)

var (
	once sync.Once
	app  *App
)

type App struct {
	cfg *cfgargs.SrvConfig
	srv *server.Server
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
	db.InitRedisClient(cfg)
	err := db.InitMongoClient(cfg)
	if err != nil {
		logger.Fatal("init mongo db err: %v", err)
		return
	}

	a.srv = server.NewServer()
	a.srv.Init(cfg)
	a.srv.Run()
}
