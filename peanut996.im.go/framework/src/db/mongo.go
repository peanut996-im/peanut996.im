package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var once sync.Once
var lastMongoClient *MongoClient

//MongoClient ..
type MongoClient struct {
	host    string
	passwd  string
	db      string
	port    string
	user    string
	session *mongo.Client
	ctx     context.Context
	timeout time.Duration
}

//GetMongoClient ...
func GetMongoClient(host, port, db, user, passwd string) *MongoClient {
	once.Do(func() {
		client, err := NewMongoClient(host, port, db, user, passwd)
		if nil != err {
			//TODO log
			// fmt.Println(err)
			panic(err)
		}
		lastMongoClient = client
	})
	return lastMongoClient
}

//NewMongoClient ...
func NewMongoClient(host, port, db, user, passwd string) (mongoClient *MongoClient, err error) {
	ctx := context.Background()

	url := fmt.Sprintf("mongodb://%v:%v@%v:%v/%v", user, passwd, host, port, db)
	option := options.Client().ApplyURI(url).SetMaxPoolSize(10)
	// s := mongo.Connect(ctx, options.Client())
	mongoSession, err := mongo.NewClient(option)
	if nil != err {
		log.Fatal("mongo connect fail. ", err)
		return nil, err
	}
	err = mongoSession.Connect(context.TODO())
	if nil != err {
		log.Fatal("mongo connect fail. ", err)
		return nil, err
	}
	mongoClient = &MongoClient{
		host:    host,
		port:    port,
		db:      db,
		passwd:  passwd,
		user:    user,
		ctx:     ctx,
		timeout: 10 * time.Second,
		session: mongoSession,
	}
	return mongoClient, nil
}

//GetAllDatabaseNames ...
func (m *MongoClient) GetAllDatabaseNames() (nameSlice []string, err error) {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	nameSlice, err = m.session.ListDatabaseNames(ctx, bson.D{})
	if nil != err {
		// TODO log
		return nil, err
	}
	return nameSlice, nil
}

//GetAllCollectionNames ...
func (m *MongoClient) GetAllCollectionNames() (collectionSlice []string, err error) {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	collectionSlice, err = m.session.Database(m.db).ListCollectionNames(ctx, bson.D{})
	if nil != err {
		// TODO log
		return nil, err
	}
	return collectionSlice, nil
}

//GetCollectionHandle ...
func (m *MongoClient) GetCollectionHandle(collection string) *mongo.Collection {
	handle := m.session.Database(m.db).Collection(collection)
	return handle
}

//Find ...
func (m *MongoClient) Find(collection string, filter interface{}, res interface{},
	opts ...*options.FindOptions) interface{} {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	handle := m.GetCollectionHandle(collection)
	cursor, err := handle.Find(ctx, filter, opts...)
	if err != nil {
		panic(err)
	}
	err = cursor.All(ctx, res)
	if nil != err {
		return err
	}
	return nil
}

//FindOne ...
func (m *MongoClient) FindOne(filter interface{},
	opts ...*options.FindOneOptions) error {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
}
