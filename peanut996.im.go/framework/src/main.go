package main

import (
	"fmt"
	"io/ioutil"

	"framework/config"

	yaml "gopkg.in/yaml.v2"
)

func init() {
	c := &config.Config{}
	configFile, err := ioutil.ReadFile("config-dev.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configFile, c)
	if nil != err {
		panic(err)
	}
	fmt.Printf("%+v\n", c)
}
func main() {
}
