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
		go s.HandleInvokeEvent(target, iR.Event, iR.Data)
	}
	logger.Info("Gate.HandleInvoke Done.")
}

func (s *Server) HandleInvokeEvent(scene, event string, data interface{}) {
	s.Lock()
	si, ok := s.SceneToSessions[scene]
	if !ok {
		logger.Info("Gate.HandleInvokeEvent Scene offline. Event: %v, Scene: %v", event, scene)
		s.Unlock()
		return
	}
	s.Unlock()
	si.Push(event, data)
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

//func (s *Server) ListenChat() {
//	if reflect.TypeOf(s.gateBroker).String() == reflect.TypeOf(&broker.GateBrokerHttp{}).String() {
//		s.gateBroker.Listen(s.ListenChatHTTP())
//	} else {
//		logger.Debug("Gate.Listen HTTP Start Failed")
//	}
//}

func (s *Server) ListenChatHTTP() interface{} {
	return func(c *gin.Context) {
		logger.Info("Gate.PushChat from Logic")
		pCR := &api.PushChatRequest{}
		err := c.BindJSON(pCR)
		if err != nil {
			logger.Error("Gate.Chat PushChat "+api.UnmarshalJsonError, err)
			c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
			return
		}
		scene := pCR.Target
		s.Lock()
		si, ok := s.SceneToSessions[scene]
		if ok {
			//online
			logger.Info("Gate.ListenChat Session Online  \n[%v]", si.ToString())
			si.Push(api.EventChat, pCR.Message)
		} else {
			//offline
			logger.Info("Gate.ListenChat Session Offline")
			messages, ok := s.offlineMessages[scene]
			if !ok {
				messages = make([]*model.ChatMessage, 0)
			}
			messages = append(messages, pCR.Message)
			s.offlineMessages[scene] = messages
			logger.Debug("%+v", s.offlineMessages[scene])
			logger.Info("Gate.ListenChat Save Message Success. ")
		}
		s.Unlock()
		c.JSON(http.StatusOK, api.NewSuccessResponse(nil))
	}
}

func (s *Server) PushOfflineMessage(session *Session) {
	//logger.Info("Gate.PushOfflineMessage Start[%v]", session.ToString())
	s.Lock()
	messages, ok := s.offlineMessages[session.GetScene()]
	if ok {
		//logger.Info("Gate.PushOfflineMessage to session[%v]", session.ToString())
		for _, message := range messages {
			session.Push(api.EventChat, message)
		}
		delete(s.offlineMessages, session.GetScene())
	} else {
		//logger.Info("Gate.PushOfflineMessage No offline Message to session[%v]", session.ToString())
	}
	s.Unlock()
	//logger.Info("Gate.PushOfflineMessage Done. Session[%v]", session.ToString())

}

func (s *Server) Debug(c *gin.Context) {
	res := "SocketIOToSessions:\n{\n"
	for socket, session := range s.SocketIOToSessions {
		res += fmt.Sprintf("    socket.id: %v, sesssion: %v\n", socket, session.ToString())
	}
	res += "}\nSceneToSessions:\n{\n"
	for scene, session := range s.SceneToSessions {
		res += fmt.Sprintf("    scene: %v, sesssion: %v\n", scene, session.ToString())
	}
	res += "}\n"
	c.String(http.StatusOK, res)
}
