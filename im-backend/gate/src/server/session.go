package server

import (
	"fmt"
	"framework/logger"
	sio "github.com/googollee/go-socket.io"
)

type Session struct {
	Conn  sio.Conn
	token string
	scene string
	query string
}

func NewSession(conn sio.Conn) *Session {
	return &Session{
		Conn:  conn,
		query: conn.URL().RawQuery,
	}
}
func (s *Session) GetScene() string {
	return s.scene
}

func (s *Session) SetScene(scene string) {
	s.scene = scene
}

func (s *Session) GetID() string {
	return s.Conn.ID()
}

func ToString(c sio.Conn) string {
	if nil != c {
		id := c.ID()
		localAddr := c.LocalAddr()
		remoteAddr := c.RemoteAddr()
		return fmt.Sprintf("ID:%v addr.local:%v addr.remote:%v", id, localAddr, remoteAddr)
	}
	return "conn not found"
}

//func (s *Session) UIDSceneString() string {
//	return fmt.Sprintf("uid:%v_scene:%v", s.uid, s.scene)
//}

func (s *Session) ToString() string {
	return fmt.Sprintf("Scene: %v, ", s.scene) + ToString(s.Conn)
}

func (s *Session) Push(event string, data interface{}) {
	logger.Debug("Gate.Push Session: [%v] Event: %v, Data: %v", s.ToString(), event, data)
	s.Conn.Emit(event, data)
}
