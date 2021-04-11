package main

import (
	"context"
	"framework/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

//User ...
type User struct {
	Name string `bson:"name"`
}

func init() {
	// c := &config.SrvConfig{}
	// configFile, err := ioutil.ReadFile(file.GetAbsPath("etc/config-example.yaml"))
	// if err != nil {
	// 	panic(err)
	// }
	// err = yaml.Unmarshal(configFile, c)
	// if nil != err {
	// 	panic(err)
	// }
	// fmt.Printf("%+v\n", c)

	// node, _ := snowflake.NewNode(1)
	// id := node.Generate()
	// fmt.Println(id.String(), id.Int64())
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://im:49@sz.peanut996.cn:27018/im"))
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://im:4pnz6V@sz.peanut996.cn:27018/im"))
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:849421294@sz.peanut996.cn:27018/"))

	// if err != nil {
	// 	//TODO log日志处理
	// 	panic(err)
	// }
	// err = client.Connect(ctx)
	// if nil != err {
	// 	panic(err)
	// }
	// // r, err := client.ListDatabases(ctx, bson.D{})
	// db, err := client.Database("im").ListCollectionNames(ctx, bson.D{})
	// if nil != err {
	// 	panic(err)
	// }
	// fmt.Println(db)
	client := db.GetMongoClient("sz.peanut996.cn", "27018", "im", "im", "4pnz6V")
	r := client.Find("user", bson.D{}, &User{})
}
