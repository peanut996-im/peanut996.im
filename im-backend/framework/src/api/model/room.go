package model

import (
	"framework/db"
	"framework/logger"
	"framework/tool"
	"go.mongodb.org/mongo-driver/bson"
)

type Room struct {
	RoomID     string `json:"roomID" bson:"roomID"`
	OneToOne   bool   `json:"oneToOne" bson:"oneToOne"`
	Status     string `json:"status" bson:"status"`
	CreateTime string `json:"createTime,omitempty" bson:"createTime"`
}

func newRoom() *Room {
	return &Room{
		RoomID:     tool.NewSnowFlakeID(),
		CreateTime: tool.GetNowUnixMilliSecond(),
		OneToOne:   false,
	}
}

func NewGroupRoom() *Room {
	r := newRoom()
	r.OneToOne = false
	return r
}

func NewFriendRoom() *Room {
	r := newRoom()
	r.OneToOne = true
	return r
}

func insertRoom(room *Room) error {
	mongo := db.GetLastMongoClient()
	res, err := mongo.InsertOne(MongoCollectionRoom, room)
	if err != nil {
		logger.Error("mongo insert room err: %v", err)
		return err
	}
	logger.Info("Mongo insert room success, id: %v", res.InsertedID)
	return nil
}

func deleteRoom(roomID string) error {
	mongo := db.GetLastMongoClient()
	filter := bson.M{"roomID": roomID}
	if _, err := mongo.DeleteOne(MongoCollectionRoom, filter); err != nil {
		return err
	}
	return nil
}

//GetRoomIDsByUID Get User-related rooms
func GetRoomIDsByUID(uid string) ([]string, error) {
	mongo := db.GetLastMongoClient()
	rooms := []string{}
	// find from db group User
	filter := bson.M{
		"uid": uid,
	}
	var (
		groupUsers []GroupUser
		friends    []Friend
	)
	err := mongo.Find(MongoCollectionGroupUser, &groupUsers, filter)
	if err != nil {
		return nil, err
	}
	for _, groupUser := range groupUsers {
		rooms = append(rooms, groupUser.GroupID)
	}

	filter = bson.M{
		"$or": bson.A{bson.M{"userA": uid}, bson.M{"userB": uid}},
	}
	err = mongo.Find(MongoCollectionFriend, &friends, filter)
	if err != nil {
		return nil, err
	}
	for _, friend := range friends {
		rooms = append(rooms, friend.RoomID)
	}
	return tool.RemoveDuplicateString(rooms), nil
}

func GetRoomsByUID(uid string) ([]*Room, error) {
	rooms := make([]*Room, 0)
	roomIDs, err := GetRoomIDsByUID(uid)
	if nil != err {
		logger.Debug("get rooms ids err: %v", err)
		return nil, err
	}
	for _, roomID := range roomIDs {
		if roomID == "" {
			logger.Warn("this is a empty room id from uid")
			continue
		}
		room, err := GetRoomByID(roomID)
		if nil != err {
			logger.Debug("get room from id. id: %v, err: %v", roomID, err)
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func GetRoomByID(roomID string) (*Room, error) {
	mongo := db.GetLastMongoClient()
	filter := bson.M{
		"roomID": roomID,
	}
	room := &Room{}
	if err := mongo.FindOne(MongoCollectionRoom, room, filter); nil != err {
		return nil, err
	}
	return room, nil
}
