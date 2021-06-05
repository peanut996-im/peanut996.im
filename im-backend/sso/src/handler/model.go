package handler

type login struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type logout struct {
	UID   string `json:"uid"`
	Token string `json:"token"`
}

type register struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func newLogin() *login {
	return &login{}
}

func newLogout() *logout {
	return &logout{}
}

func newRegister() *register {
	return &register{}
}
