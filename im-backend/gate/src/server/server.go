package server

import (
	"errors"
	"fmt"
	"framework/api"
	"framework/api/model"
	"framework/broker"
	"framework/cfgargs"
	"framework/logger"
	myhttp "framework/net/http"
	sio "github.com/googollee/go-socket.io"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	sioSrv             *sio.Server
	httpSrv            *http.Server
	offlineMessages    map[string][]*model.ChatMessage
	gateBroker         broker.GateBroker
	nsp                string
	handlers           map[string]interface{}
	SocketIOToSessions map[string]*Session
	//UIDSceneToSessions map[string]*Session
	SceneToSessions map[string]*Session
	sync.Mutex
}

func NewSIOHandlers() map[string]interface{} {
	return make(map[string]interface{})
}

func NewServer() *Server {
	s := &Server{
		sioSrv:             sio.NewServer(nil),
		nsp:                "/",
		handlers:           make(map[string]interface{}),
		offlineMessages:    make(map[string][]*model.ChatMessage),
		SocketIOToSessions: make(map[string]*Session),
		//UIDSceneToSessions: make(map[string]*Session),
		SceneToSessions: make(map[string]*Session),
	}
	return s
}

func (s *Server) MountHandlers() {
	// socket.io
	s.handlers[api.EventAddFriend] = s.SocketEventHandler(api.EventAddFriend)
	s.handlers[api.EventJoinGroup] = s.SocketEventHandler(api.EventJoinGroup)
	s.handlers[api.EventDeleteFriend] = s.SocketEventHandler(api.EventDeleteFriend)
	s.handlers[api.EventCreateGroup] = s.SocketEventHandler(api.EventCreateGroup)
	s.handlers[api.EventLeaveGroup] = s.SocketEventHandler(api.EventLeaveGroup)
	s.handlers[api.EventChat] = s.SocketEventHandler(api.EventChat)
	s.handlers[api.EventGetUserInfo] = s.SocketEventHandler(api.EventGetUserInfo)
	s.handlers[api.EventFindUser] = s.SocketEventHandler(api.EventFindUser)
	s.handlers[api.EventFindGroup] = s.SocketEventHandler(api.EventFindGroup)
	s.handlers[api.EventInviteFriend] = s.SocketEventHandler(api.EventInviteFriend)
	s.handlers[api.EventPullMessage] = s.SocketEventHandler(api.EventPullMessage)
	s.handlers[api.EventUpdateUser] = s.SocketEventHandler(api.EventUpdateUser)

	//gatebroker http handler
	path := ""
	routers := []*myhttp.Route{
		myhttp.NewRoute(api.HTTPMethodPost, "", s.HandleInvoke),
	}
	node := myhttp.NewNodeRoute(path, routers...)
	s.gateBroker.(*broker.GateBrokerHttp).AddNodeRoute(node)

	for k, v := range s.handlers {
		s.sioSrv.OnEvent(s.nsp, k, v)
	}
}

func (s *Server) Init(cfg *cfgargs.SrvConfig) {
	// rpc by http
	if cfg.Logic.Mode == "http" {
		s.gateBroker = broker.NewGateBrokerHttp()
		s.gateBroker.Init(cfg)
	}
	// sio srv init
	s.OnConnect(func(conn sio.Conn) error {
		logger.Info("socket.io connected, socket id :%v", conn.ID())
		si := NewSession(conn)
		err := s.AcceptSession(si)
		if err != nil {
			go func() {
				//Reconnect time
				switch err.Error() {
				case api.ErrorCodeToString(api.ErrorSignInvalid):
					conn.Emit("auth", api.SignInvaildResp)
				case api.ErrorCodeToString(api.ErrorTokenInvalid):
					conn.Emit("auth", api.TokenInvaildResp)
				case api.ErrorCodeToString(api.ErrorHttpInnerError):
					conn.Emit("auth", api.NewHttpInnerErrorResponse(errors.New("auth server no response.")))
				}
				<-time.After(20 * time.Second)
				conn.Close()
			}()
		} else {
			conn.Emit("auth", api.NewSuccessResponse(nil))
		}
		go func() {
			//Resend offline Message time
			time.Sleep(2 * time.Second)
			//s.PushOfflineMessage(si)
		}()
		return nil
	})

	s.OnDisconnect(func(conn sio.Conn, message string) {
		logger.Info("socket.io disconnected, socket id :%v", conn.ID())
		s.DisconnectSession(conn)
	})

	s.OnError(func(conn sio.Conn, err error) {
		logger.Error("socket.io on err: %v, id: %v", err, conn.ID())
	})

	s.MountHandlers()
	// run
	go s.Run(cfg) //nolint: errcheck
}

func (s *Server) Run(cfg *cfgargs.SrvConfig) error {
	go s.gateBroker.Listen()
	defer func(srv *sio.Server) {
		err := srv.Close()
		if err != nil {
			panic(err)
		}
	}(s.sioSrv)
	go func() {
		err := s.sioSrv.Serve()
		if err != nil {
			panic(err)
		}
	}() //nolint: errcheck

	if cfg.SocketIO.Cors {
		http.HandleFunc("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
			allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
			if origin := r.Header.Get("Origin"); origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
			}
			if r.Method == "OPTIONS" {
				return
			}
			r.Header.Del("Origin")
			s.sioSrv.ServeHTTP(w, r)
		})
	} else {
		http.Handle("/socket.io/", s.sioSrv)
	}

	addr := fmt.Sprintf(":%v", cfg.SocketIO.Port)
	logger.Info("Listening and serving Socket.IO on :%v", addr)

	err := http.ListenAndServe(addr, nil)
	logger.Fatal("Listening and serving Socket.IO at %v... err:%v", addr, err)
	return err
}

func (s *Server) OnConnect(f func(sio.Conn) error) {
	s.sioSrv.OnConnect(s.nsp, f)
}

func (s *Server) OnDisconnect(f func(sio.Conn, string)) {
	s.sioSrv.OnDisconnect(s.nsp, f)
}

func (s *Server) OnError(f func(sio.Conn, error)) {
	s.sioSrv.OnError(s.nsp, f)
}

//SetNameSpace reset namespace
func (s *Server) SetNameSpace(nsp string) {
	s.nsp = nsp
}

//AcceptSession authentication for session
func (s *Server) AcceptSession(session *Session) error {
	logger.Info("%v try to get lock")
	s.Lock()
	logger.Info("Session.Accept Start. Session[%v]", session.ToString())
	ok, err := s.Auth(session)
	if !ok || nil != err {
		s.Unlock()
		return err
	}
	s.SocketIOToSessions[session.GetID()] = session
	s.SceneToSessions[session.scene] = session
	s.Unlock()
	//logger.Info("Session.Accept succeed, session:[%v]", session.ToString())
	logger.Info("Session.Accept Done. Session[%v]", session.ToString())
	return nil
}

func (s *Server) DisconnectSession(conn sio.Conn) *Session {
	s.Lock()
	si, ok := s.SocketIOToSessions[conn.ID()]
	if ok || nil != si {
		delete(s.SocketIOToSessions, si.Conn.ID())
	} else {
		logger.Warn("Sessions.DisconnectSession[%v] not found", ToString(conn))
	}

	if nil != si {
		siScene, ok := s.SceneToSessions[si.scene]
		if ok || nil != siScene {
			logger.Info("Sessions.DisconnectSession, Scene:%v", si.scene)
			delete(s.SceneToSessions, si.scene)
		}
	}

	s.Unlock()
	return si
}
