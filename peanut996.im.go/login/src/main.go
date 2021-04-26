package main

import (
	"fmt"
	"framework/api"
	"framework/cfgargs"
	"framework/logger"
	"log"
	"net/url"
)

var (
	BuildVersion, BuildTime, BuildUser, BuildMachine string
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

	logger.InitLogger(srvConfig)

	logger.Info("App login started...")

	vals := url.Values{
		"name": []string{"Jack"},
		"age":  []string{"20"},
	}

	fmt.Println(api.MakeSign(vals, "88888888"))

}
