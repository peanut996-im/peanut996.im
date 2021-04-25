package db

import (
	"framework/cfgargs"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name string
	Age  int
}

var (
	mongoClient *MongoClient
	mongoConfig *cfgargs.SrvConfig
	user0       = User{"user0", 15}
	user1       = User{"user1", 30}
	user2       = User{"user2", 45}
)

func init() {
	cfg, err := cfgargs.InitSrvCfg(nil)
	if nil != err {
		panic("get config error")
	}
	client, err := NewMongoClient(cfg.Mongo.Host, cfg.Mongo.Port, cfg.Mongo.DB, cfg.Mongo.DB, cfg.Mongo.Passwd)
	if err != nil {
		panic(err)
	}
	mongoClient = client
	mongoConfig = cfg
}
func TestGetMongoClient(t *testing.T) {
	type args struct {
		host   string
		port   string
		db     string
		user   string
		passwd string
	}
	tests := []struct {
		name       string
		args       args
		wantClient *MongoClient
		wantErr    bool
	}{
		{"case1", args{
			host:   mongoConfig.Mongo.Host,
			port:   mongoConfig.Mongo.Port,
			db:     mongoConfig.Mongo.DB,
			user:   mongoConfig.Mongo.DB,
			passwd: mongoConfig.Mongo.Passwd}, mongoClient, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMongoClient(tt.args.host, tt.args.port, tt.args.db, tt.args.user, tt.args.passwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMongoClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMongoClient_GetAllDatabaseNames(t *testing.T) {
	tests := []struct {
		name          string
		m             *MongoClient
		wantNameSlice []string
		wantErr       bool
	}{
		{"case 0", mongoClient, []string{"im"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNameSlice, err := tt.m.GetAllDatabaseNames()
			if (err != nil) != tt.wantErr {
				t.Errorf("MongoClient.GetAllDatabaseNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNameSlice, tt.wantNameSlice) {
				t.Errorf("MongoClient.GetAllDatabaseNames() = %v, want %v", gotNameSlice, tt.wantNameSlice)
			}
		})
	}
}

func TestMongoClient_GetAllCollectionNames(t *testing.T) {
	tests := []struct {
		name                string
		m                   *MongoClient
		wantCollectionSlice []string
		wantErr             bool
	}{
		{"case0", mongoClient, []string{"user", "log"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.m.GetAllCollectionNames()
			if (err != nil) != tt.wantErr {
				t.Errorf("MongoClient.GetAllCollectionNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMongoClient_GetCollectionHandle(t *testing.T) {
	type args struct {
		collection string
	}
	tests := []struct {
		name string
		m    *MongoClient
		args args
		want *mongo.Collection
	}{
		{"case0", mongoClient, args{
			collection: "im",
		},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.GetCollectionHandle(tt.args.collection); nil == got {
				t.Errorf("MongoClient.GetCollectionHandle() return nil")
			}
		})
	}
}

func TestMongoClient_InsertOne(t *testing.T) {
	type args struct {
		collection string
		document   interface{}
		opts       []*options.InsertOneOptions
	}
	tests := []struct {
		name    string
		m       *MongoClient
		args    args
		wantErr bool
	}{
		{"case0", mongoClient, args{
			"user", user0, nil,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.m.InsertOne(tt.args.collection, tt.args.document, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MongoClient.InsertOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMongoClient_InsertMany(t *testing.T) {
	type args struct {
		collection string
		documents  []interface{}
		opts       []*options.InsertManyOptions
	}
	users := []interface{}{user1, user2}
	tests := []struct {
		name    string
		m       *MongoClient
		args    args
		wantErr bool
	}{
		{"case0", mongoClient, args{
			"user", users,
			nil,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.m.InsertMany(tt.args.collection, tt.args.documents, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MongoClient.InsertMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMongoClient_UpdateByID(t *testing.T) {
	type args struct {
		collection string
		id         interface{}
		update     interface{}
		opts       []*options.UpdateOptions
	}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.M{"author": "peanut996"}},
	}
	id, _ := primitive.ObjectIDFromHex("607cf673ef4e0dbdb044530f")
	tests := []struct {
		name    string
		m       *MongoClient
		args    args
		wantErr bool
	}{
		{"case0", mongoClient, args{
			"user", id, update, nil,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.m.UpdateByID(tt.args.collection, tt.args.id, tt.args.update, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MongoClient.UpdateByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMongoClient_UpdateOne(t *testing.T) {
	type args struct {
		collection string
		filter     interface{}
		update     interface{}
		opts       []*options.UpdateOptions
	}
	filter := bson.M{
		"name": "user2",
	}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.M{"updateBy": "peanut996"}},
	}
	tests := []struct {
		name    string
		m       *MongoClient
		args    args
		wantErr bool
	}{
		{"case0", mongoClient, args{"user", filter, update, nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.m.UpdateOne(tt.args.collection, tt.args.filter, tt.args.update, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MongoClient.UpdateOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMongoClient_UpdateMany(t *testing.T) {
	type args struct {
		collection string
		filter     interface{}
		update     interface{}
		opts       []*options.UpdateOptions
	}
	filter := bson.M{
		"name": bson.M{"$regex": primitive.Regex{
			Pattern: "^user",
			Options: "i",
		}},
	}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.M{"updateBy": "peanut996"}},
	}
	tests := []struct {
		name    string
		m       *MongoClient
		args    args
		wantErr bool
	}{
		{"case0", mongoClient, args{"user", filter, update, nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.m.UpdateMany(tt.args.collection, tt.args.filter, tt.args.update, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MongoClient.UpdateMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMongoClient_ReplaceOne(t *testing.T) {
	type args struct {
		collection  string
		filter      interface{}
		replacement interface{}
		opts        []*options.ReplaceOptions
	}
	filter := bson.M{
		"name": bson.M{"$regex": primitive.Regex{
			Pattern: "^user2$",
			Options: "i",
		}},
	}
	replacement := User{
		Name: "user2",
		Age:  100,
	}
	tests := []struct {
		name    string
		m       *MongoClient
		args    args
		wantErr bool
	}{
		{"name", mongoClient, args{"user", filter, replacement, nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.m.ReplaceOne(tt.args.collection, tt.args.filter, tt.args.replacement, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MongoClient.ReplaceOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func TestMongoClient_FindOne(t *testing.T) {
	type args struct {
		collection string
		value      interface{}
		filter     interface{}
		opts       []*options.FindOneOptions
	}
	u := User{}
	filter := bson.D{primitive.E{Key: "name", Value: "user1"}}
	tests := []struct {
		name    string
		m       *MongoClient
		args    args
		wantErr bool
	}{
		{"case0", mongoClient, args{"user", &u, filter, nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.FindOne(tt.args.collection, tt.args.value, tt.args.filter, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("MongoClient.FindOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoClient_FindOneAndDelete(t *testing.T) {
	type args struct {
		collection string
		value      interface{}
		filter     interface{}
		opts       []*options.FindOneAndDeleteOptions
	}
	u := User{}
	filter := bson.D{primitive.E{Key: "name", Value: "user2"}}
	tests := []struct {
		name    string
		m       *MongoClient
		args    args
		wantErr bool
	}{
		{"case0", mongoClient, args{"user", &u, filter, nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.FindOneAndDelete(tt.args.collection, tt.args.value, tt.args.filter, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("MongoClient.FindOneAndDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoClient_FindOneAndUpdate(t *testing.T) {
	type args struct {
		collection string
		value      interface{}
		filter     interface{}
		update     interface{}
		opts       []*options.FindOneAndUpdateOptions
	}
	u := User{}
	filter := bson.D{primitive.E{Key: "name", Value: "user0"}}
	update := bson.M{
		"$set": bson.M{
			"age": 10000,
		},
	}
	tests := []struct {
		name    string
		m       *MongoClient
		args    args
		wantErr bool
	}{
		{"case0", mongoClient, args{"user", &u, filter, update, nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.FindOneAndUpdate(tt.args.collection, tt.args.value, tt.args.filter, tt.args.update, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("MongoClient.FindOneAndUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoClient_FindOneAndReplace(t *testing.T) {
	type args struct {
		collection  string
		value       interface{}
		filter      interface{}
		replacement interface{}
		opts        []*options.FindOneAndReplaceOptions
	}
	u := User{}
	filter := bson.D{primitive.E{Key: "name", Value: "user0"}}
	replacement := User{
		Name: "user0",
		Age:  255,
	}
	tests := []struct {
		name    string
		m       *MongoClient
		args    args
		wantErr bool
	}{
		{"case0", mongoClient, args{"user", &u, filter, replacement, nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.FindOneAndReplace(tt.args.collection, tt.args.value, tt.args.filter, tt.args.replacement, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("MongoClient.FindOneAndReplace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestMongoClient_Find(t *testing.T) {
	type args struct {
		collection string
		value      interface{}
		filter     interface{}
		opts       []*options.FindOptions
	}
	users := []User{}
	filter := bson.M{
		"name": bson.M{
			"$regex": primitive.Regex{
				Pattern: "user",
				Options: "i",
			},
		},
	}
	tests := []struct {
		name string
		m    *MongoClient
		args args
	}{
		{"case0", mongoClient, args{
			collection: "im",
			value:      &users,
			filter:     filter,
			opts:       nil,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Find(tt.args.collection, tt.args.value, tt.args.filter, tt.args.opts...); nil != got {
				t.Errorf("MongoClient.Find() error: %v", got)
			}
		})
	}
}

func TestMongoClient_DeleteOne(t *testing.T) {
	type args struct {
		collection string
		filter     interface{}
		opts       []*options.DeleteOptions
	}
	filter := bson.D{primitive.E{Key: "name", Value: "user1"}}

	tests := []struct {
		name    string
		m       *MongoClient
		args    args
		wantErr bool
	}{
		{"case0", mongoClient, args{"user", filter, nil}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.m.DeleteOne(tt.args.collection, tt.args.filter, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MongoClient.DeleteOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMongoClient_DeleteMany(t *testing.T) {
	type args struct {
		collection string
		filter     interface{}
		opts       []*options.DeleteOptions
	}
	filter := bson.M{
		"name": bson.M{
			"$regex": primitive.Regex{
				Pattern: "user",
				Options: "i",
			},
		},
	}
	tests := []struct {
		name    string
		m       *MongoClient
		args    args
		wantErr bool
	}{
		{"case0", mongoClient, args{"user", filter, nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.m.DeleteMany(tt.args.collection, tt.args.filter, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MongoClient.DeleteMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
