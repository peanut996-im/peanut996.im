package model

import (
	"framework/db"
	"framework/tool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatMessage struct {
	From     string  `json:"from" bson:"from"`
	To       string  `json:"to,omitempty" bson:"to"`
	Content  string  `json:"content" bson:"content"`
	Time     string  `json:"time,omitempty" bson:"time"`
	Type     string  `json:"type" bson:"type"`
	FileName string  `json:"fileName" bson:"fileName"`
	Size     int     `json:"size" bson:"size"`
	Height   float64 `json:"height" bson:"height"`
	Width    float64 `json:"width" bson:"width"`
}

func NewChatMessage() *ChatMessage {
	return &ChatMessage{
		Time: tool.GetNowUnixNanoSecond(),
	}
}

func ChatMessageFrom(from, to, content, Type string, height, width float64, size int, fileName string) *ChatMessage {
	c := NewChatMessage()
	c.From = from
	c.To = to
	c.Content = content
	c.Type = Type
	c.Size = size
	c.Height = height
	c.Width = width
	c.FileName = fileName
	return c
}

func InsertChatMessage(c *ChatMessage) error {
	mongo := db.GetLastMongoClient()
	if _, err := mongo.InsertOne(MongoCollectionChatMessage, c); err != nil {
		return err
	}
	return nil
}

func getMessageWithPage(roomID string, current, pageSize int64) ([]*ChatMessage, error) {
	mongo := db.GetLastMongoClient()
	filter := bson.M{
		"to": roomID,
	}
	findOptions := &options.FindOptions{}
	findOptions.Sort = bson.M{"time": -1}
	findOptions.SetSkip(current)
	findOptions.SetLimit(pageSize)
	messages := make([]*ChatMessage, 0)
	err := mongo.Find(MongoCollectionChatMessage, &messages, filter, findOptions)
	if nil != err {
		return nil, err
	}
	return messages, nil
}

func GetGroupMessageWithPage(group *Group, current, pageSize int64) ([]*ChatMessage, error) {
	messages, err := getMessageWithPage(group.GroupID, current, pageSize)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// GetFriendMessageWithPage pull message with page.
// !!! Param Friend must have a roomID
func GetFriendMessageWithPage(friend *Friend, current, pageSize int64) ([]*ChatMessage, error) {
	messages, err := getMessageWithPage(friend.RoomID, current, pageSize)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
