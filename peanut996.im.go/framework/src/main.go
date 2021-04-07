package main

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

func init() {
	// c := &config.SrvConfig{}
	// configFile, err := ioutil.ReadFile(file.GetAbsPath("etc/config-example.yaml"))
	// if err != nil {
	// 	panic(err)
	// }
	// err = yaml.Unmarshal(configFile, c)
	// if nil != err {
	// 	panic(err)
	// }
	// fmt.Printf("%+v\n", c)

	node, _ := snowflake.NewNode(1)
	id := node.Generate()
	fmt.Println(id.String(), id.Int64())
}
func main() {
}
