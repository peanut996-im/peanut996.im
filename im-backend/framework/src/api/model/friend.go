// Package model
// @Title  friend.go
// @Description
// @Author  peanut996
// @Update  peanut996  2021/5/22 10:05
package model

import (
	"errors"
	"framework/db"
	"framework/logger"
	"framework/tool"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
)

type Friend struct {
	FriendA    string `json:"friendA" bson:"userA"`
	FriendB    string `json:"friendB" bson:"userB"`
	RoomID     string `json:"roomID" bson:"roomID"`
	CreateTime string `json:"-" bson:"createTime"`
}
type FriendWithRoomID struct {
	User   `json:",inline"`
	RoomID string `json:"roomID"`
}
type FriendData struct {
	User     `json:",inline"`
	RoomID   string         `json:"roomID"`
	Messages []*ChatMessage `json:"messages"`
}

func NewFriendWithRoomID(user *User, roomID string) *FriendWithRoomID {
	return &FriendWithRoomID{
		*user,
		roomID,
	}
}

func NewFriend() *Friend {
	return &Friend{
		CreateTime: tool.GetNowUnixMilliSecond(),
	}
}

func insertFriend(friend *Friend) error {
	mongo := db.GetLastMongoClient()
	r := NewFriendRoom()
	friend.RoomID = r.RoomID
	// First try to insert the room
	err := insertRoom(r)
	if nil != err {
		logger.Error("InsertFriendRoom err: %v", err)
		return err
	}
	// Second try to insert the friend relationship
	if _, err := mongo.InsertOne(MongoCollectionFriend, friend); err != nil {
		logger.Error("mongo insert friend err: %v", err)
		return err
	}
	//Symmetrical insertion
	oppositeFriend := NewFriend()
	oppositeFriend.FriendA = friend.FriendB
	oppositeFriend.FriendB = friend.FriendA
	oppositeFriend.RoomID = r.RoomID
	if _, err = mongo.InsertOne(MongoCollectionFriend, oppositeFriend); err != nil {
		logger.Error("mongo insert oppositefriend err: %v", err)
		return err
	}
	return nil
}

//AddNewFriend Add friends by UID
func AddNewFriend(friendA, friendB string) error {
	if _, err := GetFriend(friendA, friendB); nil == err {
		// already exists or error
		return errors.New("friend already exists or find error")
	}
	f := NewFriend()
	f.FriendA = friendA
	f.FriendB = friendB

	if err := insertFriend(f); nil != err {
		return err
	}
	return nil
}

func DeleteFriend(friendA, friendB string) (*Friend, error) {
	mongo := db.GetLastMongoClient()
	// find room and delete
	friend, err := GetFriend(friendA, friendB)
	if err != nil {
		return nil, err
	}
	if err := deleteRoom(friend.RoomID); nil != err {
		return nil, err
	}
	filter := bson.M{"userA": friendA, "userB": friendB}
	if _, err = mongo.DeleteMany(MongoCollectionFriend, filter); err != nil {
		return nil, err
	}
	filter = bson.M{"userB": friendA, "userA": friendB}
	if _, err = mongo.DeleteMany(MongoCollectionFriend, filter); err != nil {
		return nil, err
	}
	return friend, nil
}

func GetFriendUIDsByUID(user string) ([]string, error) {
	mongo := db.GetLastMongoClient()
	friends := make([]string, 0)
	filterA := bson.M{"userA": user}
	friendsA := []Friend{}
	if err := mongo.Find(MongoCollectionFriend, &friendsA, filterA); nil != err {
		logger.Debug("Find friendB err: %v", err)
		return nil, err
	}
	for _, friend := range friendsA {
		friends = append(friends, friend.FriendB)
	}
	filterB := bson.M{"userB": user}
	friendsB := []Friend{}
	if err := mongo.Find(MongoCollectionFriend, &friendsB, filterB); nil != err {
		logger.Debug("Find friendA err: %v", err)
		return nil, err
	}
	for _, friend := range friendsB {
		friends = append(friends, friend.FriendA)
	}
	return tool.RemoveDuplicateString(friends), nil
}

func GetFriendsByUID(uid string) ([]*Friend, error) {
	mongo := db.GetLastMongoClient()
	filter := bson.M{"userA": uid}
	var friends []*Friend
	if err := mongo.Find(MongoCollectionFriend, &friends, filter); nil != err {
		logger.Debug("Find friendB err: %v", err)
		return nil, err
	}
	return friends, nil
}

func GetFriend(friendA, friendB string) (*Friend, error) {
	mongo := db.GetLastMongoClient()
	friend := &Friend{}
	filter := bson.M{
		"userA": friendA,
		"userB": friendB,
	}

	if err := mongo.FindOne(MongoCollectionFriend, friend, filter); err != nil {
		return nil, err
	}
	return friend, nil
}

func GetFriendsByRoomID(roomID string) ([]string, error) {
	mongo := db.GetLastMongoClient()
	filter := bson.M{
		"roomID": roomID,
	}
	friends := []Friend{}
	if err := mongo.Find(MongoCollectionFriend, &friends, filter); nil != err {
		return nil, err
	}
	friendIDs := []string{}
	for _, friend := range friends {
		friendIDs = append(friendIDs, friend.FriendA)
	}
	return friendIDs, nil
}

func GetFriendWithRoomIDsByUID(uid string) ([]*FriendWithRoomID, error) {
	friendWithRoomIDs := make([]*FriendWithRoomID, 0)
	friends, err := GetFriendsByUID(uid)
	if nil != err {
		return nil, err
	}
	for _, friend := range friends {
		fR := &FriendWithRoomID{
			RoomID: friend.RoomID,
		}
		user, err := GetUserByUID(friend.FriendB)
		if err != nil {
			return nil, err
		}
		fR.User = *user
		friendWithRoomIDs = append(friendWithRoomIDs, fR)
	}
	return friendWithRoomIDs, nil
}

func GetFriendDataByFriend(friend *Friend) (*FriendData, error) {
	friendData := &FriendData{
		RoomID: friend.RoomID,
	}
	user, err := GetUserByUID(friend.FriendB)
	if err != nil {
		return nil, err
	}
	friendData.User = *user
	messages, err := GetFriendMessageWithPage(friend, 0, DefaultFriendPageSize)
	if err != nil {
		return nil, err
	}
	friendData.Messages = messages
	return friendData, nil
}

func GetFriendDataByIDs(friendA, friendB string) (*FriendData, error) {
	friend, err := GetFriend(friendA, friendB)
	if err != nil {
		return nil, err
	}
	friendData, err := GetFriendDataByFriend(friend)
	if err != nil {
		return nil, err
	}
	return friendData, nil
}

func GetFriendDatasByUID(uid string) ([]*FriendData, error) {
	friendDatas := make([]*FriendData, 0)
	friends, err := GetFriendsByUID(uid)
	if nil != err {
		return nil, err
	}
	var wg sync.WaitGroup
	var lock sync.RWMutex
	errs := make([]error, 0)
	for _, friend := range friends {
		wg.Add(1)
		go func(friend *Friend) {
			defer wg.Done()
			fD, err := GetFriendDataByFriend(friend)
			if nil != err {
				lock.Lock()
				errs = append(errs, err)
				lock.Unlock()
				return
			}
			lock.Lock()
			friendDatas = append(friendDatas, fD)
			lock.Unlock()
		}(friend)
	}
	wg.Wait()
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return friendDatas, nil
}
