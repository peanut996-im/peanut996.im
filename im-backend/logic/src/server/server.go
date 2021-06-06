package server

import (
	"framework/api"
	"framework/api/model"
	"framework/broker"
	"framework/cfgargs"
	"framework/logger"
	"framework/net/http"
	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg          *cfgargs.SrvConfig
	logicBroker  broker.LogicBroker
	httpSrv      *http.Server
	httpClient   *http.Client
	messageQueue chan *model.ChatMessage
}

func NewServer() *Server {
	return &Server{
		messageQueue: make(chan *model.ChatMessage, 5000),
	}
}

func (s *Server) Init(cfg *cfgargs.SrvConfig) {
	gin.DefaultWriter = logger.MultiWriter(logger.DefLogger().GetLogWriters()...)
	if cfg.Gate.Mode == "http" {
		s.logicBroker = broker.NewLogicBrokerHttp()
		s.logicBroker.Init(cfg)
	}
	s.logicBroker.Init(cfg)
	s.cfg = cfg
	s.httpClient = http.NewClient()
	s.httpSrv = http.NewServer()
	s.httpSrv.Init(cfg)
	s.MountRoute()

}

func (s *Server) Run() {
	go func() {
		s.Consume(s.ConsumeMessage)
	}()
	go s.logicBroker.Listen()
	//go s.httpSrv.Run()
}

func (s *Server) MountRoute() {
	path := ""
	routers := []*http.Route{
		// TODO: Mount routes
		http.NewRoute(api.HTTPMethodPost, api.EventChat, s.Chat),
		http.NewRoute(api.HTTPMethodPost, api.EventAuth, s.Auth),
		http.NewRoute(api.HTTPMethodPost, api.EventLoad, s.Load),
		http.NewRoute(api.HTTPMethodPost, api.EventAddFriend, s.AddFriend),
		http.NewRoute(api.HTTPMethodPost, api.EventDeleteFriend, s.DeleteFriend),
		http.NewRoute(api.HTTPMethodPost, api.EventCreateGroup, s.CreateGroup),
		http.NewRoute(api.HTTPMethodPost, api.EventJoinGroup, s.JoinGroup),
		http.NewRoute(api.HTTPMethodPost, api.EventLeaveGroup, s.LeaveGroup),
		http.NewRoute(api.HTTPMethodPost, api.EventGetUserInfo, s.GetUserInfo),
		http.NewRoute(api.HTTPMethodPost, api.EventFindUser, s.FindUser),
		http.NewRoute(api.HTTPMethodPost, api.EventFindGroup, s.FindGroup),
		http.NewRoute(api.HTTPMethodPost, api.EventInviteFriend, s.InviteFriend),
		http.NewRoute(api.HTTPMethodPost, api.EventPullMessage, s.PullMessage),
		http.NewRoute(api.HTTPMethodPost, api.EventUpdateUser, s.UpdateUser),
	}
	node := http.NewNodeRoute(path, routers...)
	s.logicBroker.(*broker.LogicBrokerHttp).AddNodeRoute(node)
	s.httpSrv.AddNodeRoute(node)
}

func (s *Server) Produce(message *model.ChatMessage) {
	// MQ　producer
	logger.Info("Logic.Produce: produce new message: [%+v]", *message)
	s.messageQueue <- message
}

func (s *Server) Consume(consumerFunc func(message *model.ChatMessage)) {
	for message := range s.messageQueue {
		// MQ consumer
		logger.Info("Logic.Consume: consume message: [%+v]", *message)
		//s.PushChatMessage(message)
		consumerFunc(message)
	}
}
