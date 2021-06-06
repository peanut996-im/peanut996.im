package model

import (
	"framework/db"
	"framework/logger"
	"framework/tool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

//User means a people who use the system.
type User struct {
	UID        string `json:"uid,omitempty" bson:"uid"`
	Account    string `json:"account" bson:"account"`
	Password   string `json:"-" bson:"password"`
	Avatar     string `json:"avatar,omitempty" bson:"avatar"`
	CreateTime string `json:"-" bson:"createTime"`
}

//NewUser returns a User who UID generate by snowflake Algorithm
func NewUser(account string, password string) *User {
	return &User{
		UID:        tool.NewSnowFlakeID(),
		Account:    account,
		Password:   password,
		CreateTime: tool.GetNowUnixMilliSecond(),
		Avatar:     "",
	}
}

func GetUserByAccount(account string) (*User, error) {
	mongo := db.GetLastMongoClient()
	filter := bson.M{"account": account}
	user := &User{}
	err := mongo.FindOne(MongoCollectionUser, user, filter)
	if err != nil {
		logger.Info("mongo get User from account err: %v, uid: %v", err, account)
		return nil, err
	}
	return user, nil
}

func GetUserByUID(uid string) (*User, error) {
	mongo := db.GetLastMongoClient()
	filter := bson.M{"uid": uid}
	user := &User{}
	err := mongo.FindOne(MongoCollectionUser, user, filter)
	if err != nil {
		logger.Info("mongo get User from uid err: %v, uid: %v", err, uid)
		return nil, err
	}
	return user, nil
}

func GetUsersFromUIDs(uids ...string) ([]*User, error) {
	var wg sync.WaitGroup
	var lock sync.Mutex
	var errs = make([]error, 0)
	users := make([]*User, 0)
	for _, uid := range uids {
		// bugfix: when uid is "", it will return err
		if uid == "" {
			logger.Warn("model.GetUsersFromUIDs uid is empty")
			continue
		}
		wg.Add(1)
		go func(uid string) {
			defer wg.Done()
			user, err := GetUserByUID(uid)
			if err != nil {
				lock.Lock()
				errs = append(errs, err)
				lock.Unlock()
			}
			lock.Lock()
			users = append(users, user)
			lock.Unlock()
		}(uid)
	}
	wg.Wait()
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return users, nil
}

func GetUsersByGroups(groups ...*Group) ([]*User, error) {
	uids, err := GetUserIDsByGroups(groups...)
	if nil != err {
		return nil, err
	}
	users, err := GetUsersFromUIDs(uids...)
	if nil != err {
		return nil, err
	}
	return users, nil
}

func InsertUser(user *User) error {
	mongo := db.GetLastMongoClient()
	res, err := mongo.InsertOne(MongoCollectionUser, user)
	if err != nil {
		logger.Error("mongo insert User err: %v", err)
		return err
	}
	logger.Info("Mongo insert User success, id: %v", res.InsertedID)
	return nil
}

func GetUIDFromAccount(account string) (string, error) {
	user, err := GetUserByAccount(account)
	if nil != err {
		return "", nil
	}
	return user.Account, nil
}

func FindUsersByAccount(account string) ([]*User, error) {
	mongo := db.GetLastMongoClient()
	filter := bson.M{
		"account": bson.M{
			"$regex": primitive.Regex{Pattern: ".*" + account + ".*", Options: "i"},
		},
	}
	users := make([]*User, 0)
	err := mongo.Find(MongoCollectionUser, &users, filter)
	if err != nil {
		logger.Debug("Mongo Error error: %v", err)
		return nil, err
	}
	return users, nil
}

func UpdateUser(user *User) error {
	mongo := db.GetLastMongoClient()
	filter := bson.M{
		"uid": user.UID,
	}
	_, err := mongo.ReplaceOne(MongoCollectionUser, filter, user)
	if err != nil {
		return err
	}
	return nil
}

func GetAssociatedUIDsByUID(uid string) ([]string, error) {
	uids := []string{}
	friends, err := GetFriendUIDsByUID(uid)
	if err != nil {
		return nil, err
	}
	uids = append(uids, friends...)
	groups, err := GetGroupsByUID(uid)
	if err != nil {
		return nil, err
	}
	groupUsers, err := GetUserIDsByGroups(groups...)
	if err != nil {
		return nil, err
	}
	uids = append(uids, groupUsers...)
	return tool.RemoveDuplicateString(uids), nil
}
