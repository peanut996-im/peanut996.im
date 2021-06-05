package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"framework/api"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"sync"

	"github.com/parnurzeal/gorequest"
)

var (
	once   sync.Once
	client *Client
)

//Client ...
type Client struct {
	session *http.Client
	goreq   *gorequest.SuperAgent
}

//NewClient Return a new http client.
func NewClient() *Client {
	once.Do(func() {
		client = &Client{
			goreq:   gorequest.New(),
			session: &http.Client{},
		}
	})
	return client
}

func GetAPIString(host, path string) string {
	return fmt.Sprintf("%v/%v", host, path)
}

// GetWithQueryParams Send a Get request and unmarshal.
func (c *Client) GetWithQueryParams(url string, vals url.Values, respModel interface{}) error {
	req, err := http.NewRequest(api.HTTPMethodGet, url, nil)
	if err != nil {
		return fmt.Errorf(api.NewRequestError, err, url)
	}
	req.URL.RawQuery = vals.Encode()
	resp, err := c.session.Do(req)
	if err != nil {
		return fmt.Errorf(api.DoRequestError, err, url)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf(api.ReadRespBodyError, err, url)
	}

	// respModel must be a pointer
	err = json.Unmarshal(body, respModel)
	if nil != err {
		return fmt.Errorf(api.UnmarshalJsonError, err)
	}
	return nil
}

func (c *Client) PostForm(url string, vals url.Values, respModel interface{}) error {
	resp, err := http.PostForm(url, vals)
	if nil != err {
		return fmt.Errorf(api.DoPostRequestError, err, url)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(api.ReadRespBodyError, err, url)
	}

	// respModel must be a pointer
	err = json.Unmarshal(body, respModel)
	if err != nil {
		return fmt.Errorf(api.UnmarshalJsonError, err)
	}

	return nil
}

//PostJson Post with marshal interface{} request and unmarshal response body
func (c *Client) PostJson(url string, reqModel, respModel interface{}) error {
	b, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf(api.MarshalJsonError, err)
	}

	req, err := http.NewRequest(api.HTTPMethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf(api.NewRequestError, err, url)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.session.Do(req)
	if err != nil {
		return fmt.Errorf(api.DoRequestError, err, url)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(api.ReadRespBodyError, err, url)
	}

	err = json.Unmarshal(body, respModel)
	if err != nil {
		return fmt.Errorf(api.UnmarshalJsonError, err)
	}

	return nil
}

//ObjectToUrlValues marshal interface{} into url.Values with `json` tags
func (c *Client) ObjectToUrlValues(obj interface{}) url.Values {
	vals := url.Values{}

	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Struct:
		default:
			fieldType := t.Field(i)
			val := fmt.Sprintf("%v", v.Field(i).Interface())
			if len(val) > 0 {
				vals.Add(fieldType.Tag.Get("json"), val)
			}
		}
	}

	return vals
}

func (c *Client) GetGoReq() *gorequest.SuperAgent {
	return c.goreq
}
