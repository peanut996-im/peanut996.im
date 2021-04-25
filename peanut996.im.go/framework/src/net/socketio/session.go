package socketio

type Session interface {
	Auth(string) (string, error)
	Connect()
	Disconnect()
}
