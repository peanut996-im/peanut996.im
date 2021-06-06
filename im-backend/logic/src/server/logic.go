// Package server
// @Title  logic.go
// @Description  封装业务逻辑
// @Author  peanut996
// @Update  peanut996  2021/5/31 10:43
package server

import (
	"fmt"
	"framework/api"
	"framework/api/model"
	"framework/logger"
	"framework/tool"
	"sync"
	"time"
)

func (s *Server) GetLoadData(uid string) (interface{}, error) {
	var wg sync.WaitGroup
	var lock sync.RWMutex
	user, friends, groups := &model.User{}, []*model.FriendData{}, []*model.GroupData{}
	errs := make([]error, 0)

	wg.Add(1)
	go func(uid string) {
		defer wg.Done()
		u, err := model.GetUserByUID(uid)
		if nil != err {
			lock.Lock()
			errs = append(errs, err)
			lock.Unlock()
			return
		}
		user = u
	}(uid)

	wg.Add(1)
	go func(uid string) {
		// friends
		defer wg.Done()
		fs, err := model.GetFriendDatasByUID(uid)
		if err != nil {
			lock.Lock()
			errs = append(errs, err)
			lock.Unlock()
			return
		}
		friends = fs
	}(uid)

	wg.Add(1)
	go func(uid string) {
		// group
		defer wg.Done()
		gs, err := model.GetGroupDatasByUID(uid)
		if err != nil {
			lock.Lock()
			errs = append(errs, err)
			lock.Unlock()
			return
		}
		groups = gs
	}(uid)
	wg.Wait()

	if len(errs) > 0 {
		logger.Error("Logic.LoadData err: %v", errs[0])
		return nil, errs[0]
	}
	return struct {
		User    *model.User         `json:"user"`
		Friends []*model.FriendData `json:"friends"`
		Groups  []*model.GroupData  `json:"groups"`
	}{
		user,
		friends,
		groups,
	}, nil
}

func (s *Server) PushLoadData(uid string) {
	start := time.Now()
	loadData, err := s.GetLoadData(uid)
	logger.Debug("Logic.PushLoadData /load %v", time.Since(start))
	if err != nil {
		logger.Error("Logic.PushLoadData Error: %v", err)
		return
	}
	go s.InvokeTarget(api.EventLoad, loadData, uid)
}

func (s *Server) PushChatMessage(message *model.ChatMessage) {
	roomID := message.To
	room, err := model.GetRoomByID(roomID)
	if err != nil {
		logger.Error("Logic.PushChat no such room: %v", roomID)
		return
	}
	targets := []string{}
	if room.OneToOne {
		logger.Debug("Logic.Chat Push Friend Message")
		//single
		targets, err = model.GetFriendsByRoomID(room.RoomID)
		if err != nil {
			return
		}

	} else {
		//group
		targets, err = model.GetUserIDsByGroupID(message.To)
		if err != nil {
			logger.Error("Logic.PushChat Get Group Users err: %v", err)
			return
		}
	}
	s.InvokeTarget(api.EventChat, message, targets...)
}

func (s *Server) ConsumeMessage(message *model.ChatMessage) {
	err := model.InsertChatMessage(message)
	if err != nil {
		logger.Error("Logic.ConsumeEvent err: %v", err)
	}
}

func (s *Server) InvokeTarget(event string, data interface{}, targets ...string) {
	// TODO find target on different gate nodes.
	logger.Info("Logic.InvokeTarget: event:%v, target: %v", event, targets)
	iR := &api.InvokeRequest{
		Event:   event,
		Targets: targets,
		Data:    data,
	}
	//go s.logicBroker.Invoke(iR)
	go func(data interface{}) {
		_, err := s.logicBroker.Invoke(iR)
		if err != nil {
			logger.Error("Logic.InvokeTarget Error: err: %v event:%v, target: %v, data:%v", err, event, targets, data)
			return
		}
	}(iR)
}

func (s *Server) InviteFriendsToGroup(friends []string, groupID string) error {
	for _, friend := range friends {
		err := model.CreateGroupUser(groupID, friend)
		if err != nil {
			return err
		}
	}
	// TODO 服务端主动推送客户端刷新
	return nil
}

func (s *Server) PullMessageByPage(uid, friendID, groupID string, current, pageSize int64) ([]*model.ChatMessage, error) {
	if len(friendID) > 0 {
		friend, err := model.GetFriend(uid, friendID)
		if err != nil {
			return nil, err
		}
		messages, err := model.GetFriendMessageWithPage(friend, current, pageSize)
		if err != nil {
			return nil, err
		}
		return messages, nil
	} else if len(groupID) > 0 {
		group := &model.Group{GroupID: groupID}
		messages, err := model.GetGroupMessageWithPage(group, current, pageSize)
		if err != nil {
			return nil, err
		}
		return messages, nil
	}
	return nil, api.ErrorCodeToError(api.ErrorHttpParamInvalid)
}

func (s *Server) UpdateUserInfo(uid string, account string, password string, avatar string) (*model.User, error) {
	user, err := model.GetUserByUID(uid)
	if nil != err {
		return nil, err
	}

	if len(account) > 0 {
		user.Account = account
	} else if len(password) > 0 {
		cipher := tool.EncryptBySha1(fmt.Sprintf("%v%v", password, s.cfg.AppKey))
		user.Password = cipher
	} else if len(avatar) > 0 {
		user.Avatar = avatar
	}
	err = model.UpdateUser(user)
	if nil != err {
		return nil, err
	}
	u, err := model.GetUserByUID(user.UID)
	if nil != err {
		return nil, err
	}
	return u, nil
}

func (s *Server) AddThenGetFriendData(uid, friendID string) (*model.FriendData, error) {
	err := model.AddNewFriend(uid, friendID)
	if err != nil {
		logger.Error(api.MongoDBError, err)
		return nil, err
	}
	friendData, err := model.GetFriendDataByIDs(uid, friendID)
	if err != nil {
		logger.Error(api.MongoDBError, err)
		return nil, err
	}
	return friendData, nil
}

func (s *Server) DeleteThenGetFriend(uid, friendID string) (*model.Friend, error) {
	friend, err := model.DeleteFriend(uid, friendID)
	if err != nil {
		return nil, err
	}
	return friend, nil
}

func (s *Server) CreateAndGetGroupData(groupName, groupAdmin string) (*model.GroupData, error) {
	logger.Debug("Logic.CreateAndGetGroupData Start: groupAdmin,groupName: %v,%v", groupAdmin, groupName)
	group, err := model.CreateGroup(groupName, groupAdmin)
	if err != nil {
		logger.Error(api.MongoDBError, err)
		return nil, err
	}
	// 返回数据
	groupData, err := model.GetGroupDataByGroupID(group.GroupID)
	if err != nil {
		logger.Error(api.MongoDBError, err)
		return nil, err
	}
	return groupData, nil
}

func (s *Server) JoinAndGetGroupData(uid, groupID string) (*model.GroupData, error) {
	err := model.CreateGroupUser(groupID, uid)
	if err != nil {
		logger.Error(api.MongoDBError, err)
		return nil, err
	}
	groupData, err := model.GetGroupDataByGroupID(groupID)
	if err != nil {
		logger.Error(api.MongoDBError, err)
		return nil, err
	}
	return groupData, nil
}

func (s *Server) LeaveAndGetGroupUser(uid, groupID string) (*model.GroupUser, error) {
	gUser, err := model.DeleteGroupUser(groupID, uid)
	if err != nil {
		logger.Error(api.MongoDBError, err)
		return nil, err
	}
	return gUser, nil
}
