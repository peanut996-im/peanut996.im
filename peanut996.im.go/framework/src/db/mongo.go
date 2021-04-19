package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	once            sync.Once
	lastMongoClient *MongoClient
)

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
func GetMongoClient(host, port, db, user, passwd string) (client *MongoClient, err error) {
	once.Do(func() {
		client, err = NewMongoClient(host, port, db, user, passwd)
		if nil == err {
			lastMongoClient = client
		}
	})
	return lastMongoClient, err
}

//NewMongoClient ...
func NewMongoClient(host, port, db, user, passwd string) (mongoClient *MongoClient, err error) {
	ctx := context.Background()

	url := fmt.Sprintf("mongodb://%v:%v@%v:%v/%v", user, passwd, host, port, db)
	option := options.Client().ApplyURI(url).SetMaxPoolSize(10)
	mongoSession, err := mongo.NewClient(option)
	if nil != err {
		return nil, err
	}
	err = mongoSession.Connect(context.TODO())
	if nil != err {
		return nil, err
	}
	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	err = mongoSession.Ping(newCtx, nil)
	if err != nil {
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
func (m *MongoClient) GetAllDatabaseNames() ([]string, error) {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	nameSlice, err := m.session.ListDatabaseNames(ctx, bson.D{})
	if nil != err {
		return nil, err
	}
	return nameSlice, nil
}

//GetAllCollectionNames ...
func (m *MongoClient) GetAllCollectionNames() ([]string, error) {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	collectionSlice, err := m.session.Database(m.db).ListCollectionNames(ctx, bson.D{})
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
func (m *MongoClient) Find(collection string, value, filter interface{},
	opts ...*options.FindOptions) error {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	handle := m.GetCollectionHandle(collection)
	cursor, err := handle.Find(ctx, filter, opts...)
	if err != nil {
		panic(err)
	}
	err = cursor.All(ctx, value)
	if nil != err {
		return err
	}
	return nil
}

//FindOne ...
func (m *MongoClient) FindOne(collection string, value, filter interface{},
	opts ...*options.FindOneOptions) error {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	singleResult := m.GetCollectionHandle(collection).FindOne(ctx, filter, opts...)
	return singleResult.Decode(value)
}

//FindOneAndDelete ...
func (m *MongoClient) FindOneAndDelete(collection string, value, filter interface{},
	opts ...*options.FindOneAndDeleteOptions) error {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	singleResult := m.GetCollectionHandle(collection).FindOneAndDelete(ctx, filter, opts...)
	return singleResult.Decode(value)
}

//FindOneAndUpdate ...
func (m *MongoClient) FindOneAndUpdate(collection string, value, filter,
	update interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	singleResult := m.GetCollectionHandle(collection).FindOneAndUpdate(ctx, filter, update, opts...)
	return singleResult.Decode(value)
}

//FindOneAndReplace ...
func (m *MongoClient) FindOneAndReplace(collection string, value, filter,
	replacement interface{}, opts ...*options.FindOneAndReplaceOptions) error {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	singleResult := m.GetCollectionHandle(collection).FindOneAndReplace(ctx, filter, replacement, opts...)
	return singleResult.Decode(value)
}

//InsertMany ...
func (m *MongoClient) InsertMany(collection string, documents []interface{},
	opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	return m.GetCollectionHandle(collection).InsertMany(ctx, documents, opts...)

}

//InsertOne ...
func (m *MongoClient) InsertOne(collection string, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	return m.GetCollectionHandle(collection).InsertOne(ctx, document, opts...)
}

//DeleteOne ...
func (m *MongoClient) DeleteOne(collection string, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	return m.GetCollectionHandle(collection).DeleteOne(ctx, filter, opts...)
}

//DeleteMany ...
func (m *MongoClient) DeleteMany(collection string, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	return m.GetCollectionHandle(collection).DeleteOne(ctx, filter, opts...)
}

//UpdateByID ...
func (m *MongoClient) UpdateByID(collection string, id interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	return m.GetCollectionHandle(collection).UpdateByID(ctx, id, update, opts...)
}

//UpdateOne ...
func (m *MongoClient) UpdateOne(collection string, filter, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	return m.GetCollectionHandle(collection).UpdateOne(ctx, filter, update, opts...)
}

//UpdateMany ...
func (m *MongoClient) UpdateMany(collection string, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	return m.GetCollectionHandle(collection).UpdateMany(ctx, filter, update, opts...)
}

//ReplaceOne ...
func (m *MongoClient) ReplaceOne(collection string, filter interface{},
	replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
	defer cancel()
	return m.GetCollectionHandle(collection).ReplaceOne(ctx, filter, replacement, opts...)
}
