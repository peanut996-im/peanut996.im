package socketio

import (
	"fmt"
	"framework/api"
	"framework/cfgargs"
	"framework/logger"
	"net/http"
	"net/url"
	"sync"

	sio "github.com/googollee/go-socket.io"
)

type Server struct {
	srv                *sio.Server
	nsp                string
	handlers           map[string]interface{}
	SocketIOToSessions map[string]*Session
	UIDSceneToSessions map[string]*Session
	UIDToSessions      map[string]*Session
	sync.Mutex
}

func NewSIOHandlers() map[string]interface{} {
	return make(map[string]interface{})
}

func NewServer() *Server {
	s := &Server{
		srv:                sio.NewServer(nil),
		nsp:                "/",
		handlers:           make(map[string]interface{}),
		SocketIOToSessions: make(map[string]*Session),
		UIDSceneToSessions: make(map[string]*Session),
		UIDToSessions:      make(map[string]*Session),
	}
	return s
}

func (s *Server) AcceptSession(session *Session, query string) (error int) {
	vals, err := url.ParseQuery(query)
	sign, _ := api.MakeSignWithQueryParams(vals, cfgargs.GetLastSrvConfig().AppKey)
	if sign != vals.Get("sign") {
		logger.Info("Session[%v]'s  sign: %v", session.ToString(), sign)
		return api.ErrorSignInvalid
	}
	if nil != err {
		logger.Info("parse token failed, err: %v", err)
		return api.ErrorTokenInvalid
	}

	t := vals.Get("token")
	_, err = session.Auth(t)
	if err != nil {
		logger.Info("token not valid, session:[%v], err:[%v]", session.ToString(), error)
		return api.ErrorTokenInvalid
	}
	logger.Info("Session.Accept succeed, session:[%v]", session.ToString())

	s.Lock()
	s.SocketIOToSessions[session.sid] = session
	s.UIDToSessions[session.uid] = session
	s.Unlock()

	logger.Info("Session.Accept done. Session[%v]", session.ToString())
	return api.ErrorCodeOK

}

func (s *Server) Run(cfg *cfgargs.SrvConfig) error {

	defer func(srv *sio.Server) {
		err := srv.Close()
		if err != nil {
			panic(err)
		}
	}(s.srv)
	go func() {
		err := s.srv.Serve()
		if err != nil {
			panic(err)
		}
	}() //nolint: errcheck

	if cfg.HTTP.Cors {
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
			s.srv.ServeHTTP(w, r)
		})
	} else {
		http.Handle("/socket.io/", s.srv)
	}

	addr := fmt.Sprintf(":%v", cfg.SocketIO.Port)
	logger.Info("Serving at %v...", addr)

	err := http.ListenAndServe(addr, nil)
	logger.Fatal("Serving at %v... err:%v", addr, err)
	return err
}

func (s *Server) OnConnect(f func(sio.Conn) error) {
	s.srv.OnConnect(s.nsp, f)
}

func (s *Server) OnDisconnect(f func(sio.Conn, string)) {
	s.srv.OnDisconnect(s.nsp, f)
}

func (s *Server) OnError(f func(sio.Conn, error)) {
	s.srv.OnError(s.nsp, f)
}

func (s *Server) MountHandlers(nsp string, handlers map[string]interface{}) {
	for k, v := range handlers {
		//fmt.Printf("nsp is %v",nsp)
		s.srv.OnEvent(nsp, k, v)
	}
}

func (s *Server) SocketIOToSession(c sio.Conn) *Session {
	s.Lock()
	si, ok := s.SocketIOToSessions[c.ID()]
	s.Unlock()
	if !ok {
		logger.Warn("session not found")
		return nil
	}
	return si
}

func (s *Server) UIDSceneToSession(uidScene string) *Session {
	s.Lock()
	si, ok := s.UIDSceneToSessions[uidScene]
	s.Unlock()
	if !ok {
		logger.Warn("session not found")
		return nil
	}
	return si
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
		siScene, ok := s.UIDSceneToSessions[si.UIDSceneString()]
		if ok || nil != siScene {
			logger.Info("Sessions.DisconnectSession,UIDAndScene:v%", si.UIDSceneString())
			delete(s.UIDSceneToSessions, si.UIDSceneString())
		}
	}

	s.Unlock()
	return si
}

//SetNameSpace 改变默认的namespace
func (s *Server) SetNameSpace(nsp string) {
	s.nsp = nsp
}
