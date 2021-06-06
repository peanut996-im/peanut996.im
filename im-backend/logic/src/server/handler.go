package server

import (
	"framework/api"
	"framework/api/model"
	"framework/db"
	"framework/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Chat 推送信息
func (s *Server) Chat(c *gin.Context) {
	cR := &api.ChatRequest{}
	err := c.BindJSON(cR)
	if err != nil {
		logger.Error("Logic.Auth "+api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	msg := model.ChatMessageFrom(cR.From, cR.To, cR.Content, cR.Type, cR.Height, cR.Width, cR.Size, cR.FileName)

	go s.PushChatMessage(msg)
	go s.Produce(msg)
	c.JSON(http.StatusOK, api.NewSuccessResponse(nil))
}

// GetUserInfo 获取用户信息
func (s *Server) GetUserInfo(c *gin.Context) {
	uR := &api.UserRequest{}
	err := c.BindJSON(uR)
	if err != nil {
		logger.Error("Logic.Auth "+api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	user, err := model.GetUserByUID(uR.UID)
	if err != nil {
		logger.Error("Logic.GetUserInfo "+api.MongoDBError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, api.NewSuccessResponse(user))
}

// Auth 用户鉴权 鉴权成功后推送初始化信息
func (s *Server) Auth(c *gin.Context) {
	aR := &api.AuthRequest{}
	err := c.BindJSON(aR)
	if err != nil {
		logger.Error("Logic.Auth "+api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	user, err := api.CheckToken(aR.Token)
	if err != nil {
		if db.IsNotExistError(err) {
			// token expired
			c.AbortWithStatusJSON(http.StatusOK, api.TokenInvaildResp)
			return
		}
		logger.Error(api.MongoDBError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	defer func(uid string) {
		// Auth success then push load data
		logger.Debug("Logic.Auth defer. uid: %v", uid)
		go s.PushLoadData(uid)
	}(user.UID)
	c.JSON(http.StatusOK, api.NewSuccessResponse(user))
}

// Load 推送初始化信息
func (s *Server) Load(c *gin.Context) {
	lR := &api.LoadRequest{}
	err := c.BindJSON(lR)
	if err != nil {
		logger.Error(api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	go s.PushLoadData(lR.UID)
	c.JSON(http.StatusOK, api.NewSuccessResponse(nil))
}

// AddFriend 添加好友
func (s *Server) AddFriend(c *gin.Context) {
	fR := &api.FriendRequest{}
	err := c.BindJSON(fR)
	if err != nil {
		logger.Error(api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	friendData, err := s.AddThenGetFriendData(fR.FriendA, fR.FriendB)
	if err != nil {
		logger.Error("Logic.AddFriend failed. err: %v", err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	defer func() {
		go s.PushLoadData(fR.FriendA)
		go s.PushLoadData(fR.FriendB)
	}()
	c.JSON(http.StatusOK, api.NewSuccessResponse(friendData))
}

// DeleteFriend 删除好友
func (s *Server) DeleteFriend(c *gin.Context) {
	fR := &api.FriendRequest{}
	err := c.BindJSON(fR)
	if err != nil {
		logger.Error(api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	friend, err := s.DeleteThenGetFriend(fR.FriendA, fR.FriendB)
	if err != nil {
		logger.Error(api.MongoDBError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	defer func() {
		go s.PushLoadData(fR.FriendA)
		go s.PushLoadData(fR.FriendB)
	}()
	c.JSON(http.StatusOK, api.NewSuccessResponse(friend))
}

// CreateGroup 创建群组
func (s *Server) CreateGroup(c *gin.Context) {
	gR := &api.GroupRequest{}
	err := c.BindJSON(gR)
	if err != nil {
		logger.Error(api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	logger.Debug("Logic.CreateGroup group name: %v uid: %v ", gR.GroupName, gR.UID)
	groupData, err := s.CreateAndGetGroupData(gR.GroupName, gR.UID)
	if err != nil {
		logger.Error(api.MongoDBError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, api.NewSuccessResponse(groupData))
}

// JoinGroup 加入群组并广播
func (s *Server) JoinGroup(c *gin.Context) {
	gR := &api.GroupRequest{}
	err := c.BindJSON(gR)
	if err != nil {
		logger.Error(api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	groupData, err := s.JoinAndGetGroupData(gR.UID, gR.GroupID)
	if err != nil {
		logger.Error(api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	defer func() {
		for _, member := range groupData.Members {
			go s.PushLoadData(member.UID)
		}
	}()
	c.JSON(http.StatusOK, api.NewSuccessResponse(groupData))
}

// LeaveGroup 离开群组并广播
func (s *Server) LeaveGroup(c *gin.Context) {
	gR := &api.GroupRequest{}
	err := c.BindJSON(gR)
	if err != nil {
		logger.Error(api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	gUser, err := s.LeaveAndGetGroupUser(gR.UID, gR.GroupID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	defer func(groupID string) {
		uids, err := model.GetUserIDsByGroupID(groupID)
		if err != nil {
			logger.Error("Logic.PushLoadData After LeaveGroup Failed. err: %v", err)
			return
		}
		for _, uid := range uids {
			go s.PushLoadData(uid)
		}
	}(gR.GroupID)
	c.JSON(http.StatusOK, api.NewSuccessResponse(gUser))
}

// FindUser 模糊搜索用户
func (s *Server) FindUser(c *gin.Context) {
	fUR := &api.FindRequest{}
	err := c.BindJSON(fUR)
	if nil != err {
		logger.Error(api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	users, err := model.FindUsersByAccount(fUR.Account)
	if nil != err {
		logger.Error(api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, api.NewSuccessResponse(users))
}

// FindGroup 模糊搜索群组
func (s *Server) FindGroup(c *gin.Context) {
	fUR := &api.FindRequest{}
	err := c.BindJSON(fUR)
	if nil != err {
		logger.Error(api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	groups, err := model.FindGroupsByGroupName(fUR.GroupName)
	if nil != err {
		logger.Error(api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, api.NewSuccessResponse(groups))
}

// InviteFriend 邀请好友进群
func (s *Server) InviteFriend(c *gin.Context) {
	iR := &api.InviteRequest{}
	err := c.BindJSON(iR)
	if err != nil {
		logger.Error(api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	err = s.InviteFriendsToGroup(iR.Friends, iR.GroupID)
	if err != nil {
		logger.Error("InviteFriendsToGroup err: %v", err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	defer func(groupID string) {
		uids, err := model.GetUserIDsByGroupID(groupID)
		if err != nil {
			logger.Error("Logic.PushLoadData After LeaveGroup Failed. err: %v", err)
			return
		}
		for _, uid := range uids {
			go s.PushLoadData(uid)
		}
	}(iR.GroupID)
	c.JSON(http.StatusOK, api.NewSuccessResponse(nil))
}

// PullMessage 分页加载信息
func (s *Server) PullMessage(c *gin.Context) {
	pR := &api.PullRequest{}
	err := c.BindJSON(pR)
	if err != nil {
		logger.Error(api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	messages, err := s.PullMessageByPage(pR.UID, pR.FriendID, pR.GroupID, pR.Current, pR.PageSize)
	if err != nil {
		logger.Error("PullMessage err: %v", err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, api.NewSuccessResponse(messages))
}

// UpdateUser 更新用户信息
func (s *Server) UpdateUser(c *gin.Context) {
	uR := &api.UpdateRequest{}
	err := c.BindJSON(uR)
	if err != nil {
		logger.Error(api.UnmarshalJsonError, err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	user, err := s.UpdateUserInfo(uR.UID, uR.Account, uR.Password, uR.Avatar)
	if nil != err {
		logger.Error("UpdateUserInfo err: %v", err)
		c.AbortWithStatusJSON(http.StatusOK, api.NewHttpInnerErrorResponse(err))
		return
	}
	defer func(uid string) {
		targets, err := model.GetAssociatedUIDsByUID(uid)
		if err != nil {
			return
		}
		for _, target := range targets {
			go s.PushLoadData(target)
		}
	}(user.UID)
	c.JSON(http.StatusOK, api.NewSuccessResponse(user))
}
