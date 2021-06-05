// Package model
// @Title  groupuser.go
// @Description
// @Author  peanut996
// @Update  peanut996  2021/5/22 17:40
package model

import (
	"framework/db"
	"framework/logger"
	"framework/tool"
	"go.mongodb.org/mongo-driver/bson"
)

type GroupUser struct {
	GroupID    string `json:"groupID" bson:"groupID"`
	UID        string `json:"uid" bson:"uid"`
	CreateTime string `json:"-" bson:"createTime"`
}

func NewGroupUser() *GroupUser {
	return &GroupUser{
		CreateTime: tool.GetNowUnixMilliSecond(),
	}
}

func insertGroupUser(groupUser *GroupUser) error {
	mongo := db.GetLastMongoClient()
	_, err := mongo.InsertOne(MongoCollectionGroupUser, groupUser)
	if err != nil {
		logger.Error("mongo insert GroupUser err: %v", err)
		return err
	}
	return nil
}

func CreateGroupUser(groupID, user string) error {
	gp := NewGroupUser()
	gp.GroupID = groupID
	gp.UID = user
	if err := insertGroupUser(gp); nil != err {
		logger.Error("mongo insert GroupUser err: %v", err)
		return err
	}
	return nil
}

func DeleteGroupUser(groupID, user string) (*GroupUser, error) {
	mongo := db.GetLastMongoClient()
	filter := bson.M{"groupID": groupID, "uid": user}
	gUser := &GroupUser{}
	if err := mongo.FindOneAndDelete(MongoCollectionGroupUser, gUser, filter); nil != err {
		return nil, err
	}
	return gUser, nil
}

func GetUserIDsByGroupID(groupID string) ([]string, error) {
	mongo := db.GetLastMongoClient()
	filter := bson.M{"groupID": groupID}
	gUsers := make([]GroupUser, 0)
	userIDs := make([]string, 0)
	if err := mongo.Find(MongoCollectionGroupUser, &gUsers, filter); nil != err {
		return nil, err
	}
	for _, user := range gUsers {
		userIDs = append(userIDs, user.UID)
	}
	return userIDs, nil
}

func GetGroupIDsByUID(uid string) ([]string, error) {
	mongo := db.GetLastMongoClient()
	filter := bson.M{"uid": uid}
	gUsers := make([]GroupUser, 0)
	if err := mongo.Find(MongoCollectionGroupUser, &gUsers, filter); nil != err {
		return nil, err
	}
	groupIDs := make([]string, 0)
	for _, groupUser := range gUsers {
		groupIDs = append(groupIDs, groupUser.GroupID)
	}
	return groupIDs, nil
}

func GetUserIDsByGroup(group *Group) ([]string, error) {
	uids, err := GetUserIDsByGroupID(group.GroupID)
	if nil != err {
		return nil, err
	}
	return uids, nil
}

func GetUserIDsByGroups(groups ...*Group) ([]string, error) {
	uids := []string{}
	for _, group := range groups {
		us, err := GetUserIDsByGroup(group)
		if err != nil {
			return nil, err
		}
		uids = append(uids, us...)
	}
	return uids, nil
}

func GetUsersByGroup(group *Group) ([]*User, error) {
	uids, err := GetUserIDsByGroup(group)
	if nil != err {
		return nil, err
	}
	users, err := GetUsersFromUIDs(uids...)
	if err != nil {
		return nil, err
	}
	return users, nil
}
