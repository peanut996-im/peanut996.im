// Package server
// @Title  handler.go
// @Description
// @Author  peanut996
// @Update  peanut996  2021/5/24 16:11
package server

import (
	"encoding/json"
	"fmt"
	"framework/api"
	"framework/api/model"
	"framework/cfgargs"
	"framework/logger"
	"framework/tool"
	"github.com/gin-gonic/gin"
	sio "github.com/googollee/go-socket.io"
	"net/http"
	"net/url"
	"strings"
)

func (s *Server) HandleInvoke(c *gin.Context) {
	logger.Info("Gate.HandleInvoke from Logic")
	iR := &api.InvokeRequest{}
	err := c.BindJSON(iR)
	if err != nil {
		logger.Error("Gate.HandleInvoke "+api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	for _, target := range iR.Targets {
		s.Produce(&api.SingleInvokeRequest{Target: target, Event: iR.Event, Data: iR.Data})
	}
	logger.Info("Gate.HandleInvoke Done.")
}

func (s *Server) HandleInvokeEvent(scene, event string, data interface{}) {
	// TODO multiply lock race.
	s.Lock()
	sessions, ok := s.SceneToSessions[scene]
	if !ok {
		logger.Info("Gate.HandleInvokeEvent Scene offline. Event: %v, Scene: %v", event, scene)
		s.Unlock()
		return
	}
	s.Unlock()
	for _, si := range sessions {
		go si.Push(event, data)
	}

}

func (s *Server) SocketEventHandler(event string) interface{} {
	return func(conn sio.Conn, data interface{}) {
		logger.Info("/%v from[%v]: %+v", event, conn.ID(), data)
		rawJson, err := s.gateBroker.Send(event, data)
		if nil != err {
			conn.Emit(event, api.NewHttpInnerErrorResponse(err))
			logger.Error("Gate.Event[%v] Broker err: %v", event, err)
		}
		conn.Emit(event, rawJson.(json.RawMessage))
	}
}

func (s *Server) Auth(session *Session) (bool, error) {
	vals, err := url.ParseQuery(session.query)
	sign, _ := api.MakeSignWithQueryParams(vals, cfgargs.GetLastSrvConfig().AppKey)
	if sign != vals.Get("sign") {
		logger.Info("Session.Auth failed. sign invalid: %v", sign)
		return false, api.ErrorCodeToError(api.ErrorSignInvalid)
	}
	if nil != err {
		logger.Info("parse token failed, err: %v", err)
		return false, api.ErrorCodeToError(api.ErrorTokenInvalid)
	}

	t := vals.Get("token")
	rawJson, err := s.gateBroker.Send(api.EventAuth, t)
	if err != nil {
		logger.Error("Session.Auth get auth response err. err: %v", err)
		return false, api.ErrorCodeToError(api.ErrorHttpInnerError)
	}
	resp := &api.BaseRepsonse{}
	if err = json.Unmarshal(rawJson.(json.RawMessage), resp); err != nil {
		logger.Info("Session.Save json unmarshal err. err:%v, Session:[%v]", err, session.ToString())
		return false, api.ErrorCodeToError(api.ErrorHttpInnerError)
	}
	if resp.Code != api.ErrorCodeOK || resp.Data == nil {
		// Auth failed
		logger.Error("Session.Auth auth failed. Maybe token expired or user not exist? UID: [%v], Session:[%v]", vals.Get("uid"), session.ToString())
		return false, api.ErrorCodeToError(api.ErrorTokenInvalid)
	}
	u := &model.User{}
	if err = tool.MapToStruct(resp.Data, u); err != nil {
		logger.Info("Session.Auth json unmarshal err. err:%v, Session:[%v]", err, session.ToString())
		return false, api.ErrorCodeToError(api.ErrorHttpInnerError)
	}

	//logger.Info("Session.Auth succeed.")
	session.SetScene(u.UID)
	session.token = t
	return true, nil
}

func (s *Server) ConsumeEvent(event *api.SingleInvokeRequest) {
	s.HandleInvokeEvent(event.Target, event.Event, event.Data)
}

func (s *Server) DebugMapVars(c *gin.Context) {
	sb := strings.Builder{}
	sb.WriteString("SceneToSessions: \n{\n\n")
	//SceneToSessionsString,_ := tool.PrettyPrint(s.SceneToSessions)
	for key, sessions := range s.SceneToSessions {
		tmp := strings.Builder{}
		tmp.WriteString("[")
		for _, session := range sessions {
			tmp.WriteString(fmt.Sprintf("[%v],", session.ToString()))
		}
		tmp.WriteString("]\n")
		sb.WriteString(fmt.Sprintf("  %v: %v\n", key, tmp.String()))
	}
	sb.WriteString("}\nSocketIOToSessions: \n{\n\n")
	for key, session := range s.SocketIOToSessions {
		sb.WriteString(fmt.Sprintf("  %v: %v\n\n", key, session.ToString()))
	}
	sb.WriteString("}\n")
	c.String(http.StatusOK, sb.String())
}
