// Package gate
// @Title  logicbroker.go
// @Description
// @Author  peanut996
// @Update  peanut996  2021/5/31 11:39
package broker

import "framework/cfgargs"

type LogicBroker interface {
	Init(*cfgargs.SrvConfig)
	Listen()
	Invoke(packet interface{}) (interface{}, error)
	InvokeTarget(target, event string, data interface{}) (interface{}, error)
}
