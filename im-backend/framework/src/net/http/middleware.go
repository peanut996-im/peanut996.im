package http

import (
	"bytes"
	"framework/api"
	"framework/cfgargs"
	"framework/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

// middlewares must be use in app service instead of framework.
// cors should be front of checksign.

func CheckSign(cfg *cfgargs.SrvConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		b, err := c.GetRawData()
		if err != nil {
			logger.Error("get raw data err: %v", err)
			c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
			return
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))

		body := string(b)

		if err != nil || len(body) == 0 {
			// query string
			logger.Debug("get http query param: %v", c.Request.URL.Query())
			checkResult, err := api.CheckSignFromQueryParams(c.Request.URL.Query(), cfg.AppKey)
			if !checkResult || err != nil {
				logger.Debug("check sign with query string failed: query params: %v", c.Request.URL.Query())
				if !cfg.HTTP.Release {
					sign, err := api.MakeSignWithQueryParams(c.Request.URL.Query(), cfg.AppKey)
					if err == nil {
						c.AbortWithStatusJSON(http.StatusOK, api.NewBaseResponse(api.ErrorSignInvalid, gin.H{"sign": sign}))
						return
					}
				}
				c.AbortWithStatusJSON(http.StatusOK, api.SignInvaildResp)
				return
			}
		} else {
			// json
			logger.Debug("get http body: %v", body)
			checkResult, err := api.CheckSignFromJsonParams(body, cfg.AppKey)
			if !checkResult || err != nil {
				logger.Debug("check sign with json failed: body: %v", body)
				if !cfg.HTTP.Release {
					sign, err := api.MakeSignWithJsonParams(body, cfg.AppKey)
					if err == nil {
						c.AbortWithStatusJSON(http.StatusOK, api.NewBaseResponse(api.ErrorSignInvalid, gin.H{"sign": sign}))
						return
					}
				}
				c.AbortWithStatusJSON(http.StatusOK, api.SignInvaildResp)
				return
			}
		}
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return cors.Default()
}

func IPWhiteList(whitelist map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !whitelist[ip] {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
