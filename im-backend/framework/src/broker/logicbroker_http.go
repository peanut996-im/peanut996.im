// Package gate
// @Title  logicbroker_http.go
// @Description
// @Author  peanut996
// @Update  peanut996  2021/5/31 11:40
package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"framework/cfgargs"
	"framework/logger"
	"framework/net/http"
)

type LogicBrokerHttp struct {
	srv      *http.Server
	client   *http.Client
	gateAddr string
}

func NewLogicBrokerHttp() *LogicBrokerHttp {
	return &LogicBrokerHttp{}
}

func (l *LogicBrokerHttp) Init(cfg *cfgargs.SrvConfig) {
	l.srv = http.NewServer()
	l.srv.Init(cfg)
	l.client = http.NewClient()
	l.gateAddr = fmt.Sprintf("http://%v:%v", cfg.Gate.Host, cfg.Gate.Port)
	if cfg.Gate.Mode != "http" {
		logger.Warn("can't load gate configuration for http")
		if cfg.Gate.Panic {
			panic(errors.New("can't load gate configuration for http"))
		}
		return
	}
}

func (l *LogicBrokerHttp) Listen() {
	//启动http server
	l.srv.Run()
}

//Invoke Push data to gate
func (l *LogicBrokerHttp) Invoke(packet interface{}) (interface{}, error) {
	// FIXME With target, event
	// default http address
	addr := l.gateAddr + "/"
	resp, body, errs := l.client.GetGoReq().Post(addr).Send(packet).End()
	if len(errs) != 0 {
		for i, err := range errs {
			logger.Info("LogicBroker Event /%v failed. errs[%v]: %v ", i, err)
		}
		return nil, errs[0]
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("LogicBroker Event /%v http failed code: %v", resp.StatusCode))
	}
	return json.RawMessage(body), nil
}

// InvokeTarget Send commands to the gate instance
func (l *LogicBrokerHttp) InvokeTarget(target, event string, data interface{}) (interface{}, error) {
	panic("unimplemented")
}

// AddNodeRoute Mount the route to the internal HTTP server.
func (l *LogicBrokerHttp) AddNodeRoute(nodes ...*http.NodeRoute) {
	l.srv.AddNodeRoute(nodes...)
}
