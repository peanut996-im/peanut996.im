package socketio

import (
	"framework/cfgargs"
	"net/http"

	sio "github.com/googollee/go-socket.io"
)

type Server interface {
	AcceptSession(Session) error
	SceneSession(sio.Conn, string) Session
	SocketIOToSession(sio.Conn) Session
	SessionToSocketIO(Session) sio.Conn
	SceneToSession(string) Session
	SIOHandler(srvCfg *cfgargs.SrvConfig) http.HandlerFunc
	MountHandlers(map[string]func(sio.Conn, ...interface{}) interface{})
}
