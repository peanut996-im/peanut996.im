package http

import (
	"net/http"

	"github.com/parnurzeal/gorequest"
)

//Client ...
type Client struct {
	session *http.Client
	client  *gorequest.SuperAgent
}
