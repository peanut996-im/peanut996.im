package config

import (
	"fmt"
	"framework/file"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

//SrvConfig records for all conf
type SrvConfig struct {
	Mongo `yaml:"mongo"`
	Redis `yaml:"redis"`
}

//Mongo conf for mongoDB
type Mongo struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	DB     string `yaml:"db"`
	Passwd string `yaml:"passwd"`
}

//Redis configure for Redis
type Redis struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Passwd string `yaml:"passwd"`
	DB     int    `yaml:"db"`
}

//GetSrvConfig return a SrvConfig.
func GetSrvConfig(filePath string) (config *SrvConfig, err error) {
	configFile, err := ioutil.ReadFile(file.GetAbsPath(filePath))
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(configFile, &config)
	if nil != err {
		log.Fatal("Get Server SrvConfig Failed. Error: ", err)
	}
	return config, err
}

// GetRedisAddr printf the redis addr from srvconfig
func GetRedisAddr(config *SrvConfig) string {
	return fmt.Sprintf("%v:%v", config.Redis.Host, config.Redis.Port)
}
