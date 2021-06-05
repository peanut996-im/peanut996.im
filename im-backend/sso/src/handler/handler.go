package handler

import (
	"fmt"
	"framework/api/model"
	"framework/tool"
	"net/http"

	"framework/api"
	"framework/cfgargs"
	"framework/db"
	"framework/logger"
	"github.com/gin-gonic/gin"
)

//SignIn 用户登录
func SignIn(c *gin.Context) {
	login := newLogin()
	err := c.ShouldBind(login)
	if err != nil {
		c.JSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	logger.Debug("mongo find account: %v", login.Account)

	user, err := model.GetUserByAccount(login.Account)
	if err != nil {
		// user not found
		logger.Info("mongo user not found. err:%v", err)
		c.JSON(http.StatusOK, api.ResourceNotFoundResp)
		return
	}

	if user.Password == tool.EncryptBySha1(fmt.Sprintf("%v%v", login.Password, cfgargs.GetLastSrvConfig().AppKey)) {
		token, err := api.InsertToken(user.UID)
		if err != nil {
			// token failed
			logger.Error("insert user token failed")
			c.JSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
			return
		}
		// login success
		c.JSON(http.StatusOK, api.NewSuccessResponse(gin.H{
			"uid":   user.UID,
			"token": token,
		}))
		return
	}

	logger.Info("user passwd not correct")
	c.JSON(http.StatusOK, api.AuthFaildResp)
}

func SignOut(c *gin.Context) {
	logout := newLogout()

	err := c.ShouldBind(logout)
	if err != nil {
		c.JSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}

	_, err = api.CheckToken(logout.Token)
	if err != nil {
		if db.IsNotExistError(err) {
			c.JSON(http.StatusOK, api.TokenInvaildResp)
			return
		}
		c.JSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}

	err = api.DeleteToken(logout.Token)
	if err != nil {
		c.JSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, api.NewSuccessResponse(nil))
}

//SignUp 用户注册
func SignUp(c *gin.Context) {
	register := newRegister()

	err := c.ShouldBind(register)
	if err != nil {
		c.JSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	// 先看是否重复
	_, err = model.GetUserByAccount(register.Account)
	if db.IsNoDocumentError(err) {
		user := model.NewUser(register.Account, tool.EncryptBySha1(fmt.Sprintf("%v%v", register.Password, cfgargs.GetLastSrvConfig().AppKey)))
		if err = model.InsertUser(user); err != nil {
			logger.Error("SSO.Register err: %v", err)
			c.JSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
			return
		}

		// 插入token
		token, err := api.InsertToken(user.UID)
		if err != nil {
			logger.Error("redis insert user token err: %v", err)
			c.JSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
			return
		}

		c.JSON(200, api.NewSuccessResponse(gin.H{
			"uid":   user.UID,
			"token": token,
		}))

	} else {
		logger.Info("mongo user already exists or err: %v", err)
		c.JSON(http.StatusOK, api.ResourceExistsResp)
	}

}
