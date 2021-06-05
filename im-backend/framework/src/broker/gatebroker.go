package broker

import "framework/cfgargs"

type GateBroker interface {
	Init(cfg *cfgargs.SrvConfig)
	Send(string, interface{}) (interface{}, error)
	Listen()
	Register()
}
