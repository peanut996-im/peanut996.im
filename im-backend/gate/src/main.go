package main

import (
	"framework/cfgargs"
	"framework/logger"
	"gate/app"
	"log"
	"os"
	"os/signal"
)

var (
	BuildTime    string
	BuildUser    string
	BuildMachine string
	BuildVersion string
)

func main() {
	build := &cfgargs.Build{
		BuildVersion: BuildVersion,
		BuildTime:    BuildTime,
		BuildUser:    BuildUser,
		BuildMachine: BuildMachine,
	}
	srvConfig, err := cfgargs.InitSrvCfg(build, nil)
	if err != nil {
		log.Fatal(err)
	}
	srvConfig.Print()
	logger.InitLogger(srvConfig)
	app.GetApp().Init(srvConfig)
	logger.Info("App gate started...")

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Fatal("Server force shutdown...")
	close(quit)
	os.Exit(1)
}
