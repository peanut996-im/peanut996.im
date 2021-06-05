# Bug

记录写代码遇到的一些Bug

### Go

* 数据库ID选用SnowFlake算法生成，但是注意最好不要使用int64来接收，不仅仅前端，以及数据库查询时也会容易忘记转换。
* Redis和Mongo时刻记得判断Nil值，分别是`Redis.Nil` 和 `Mongo.ErrNoDocuments`



