package db

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var lastRdsClient *RedisClient

type RedisClient struct {
	addr              string
	pass              string
	db                int
	session           *redis.Client
	keepAliveInterval time.Duration // second
}

func IsNotExistError(err error) bool {
	if err == redis.Nil {
		return true
	}
	return false
}

func GetLastRedisClient() *RedisClient {
	return lastRdsClient
}

func NewRedisClient(addr string, pass string, db int, panicIfDisconnect bool) *RedisClient {
	redisClient := &RedisClient{
		addr: addr,
		pass: pass,
		db:   db,
		session: redis.NewClient(&redis.Options{
			Addr:               addr,
			Password:           pass,
			DB:                 db,
			PoolSize:           0xFF,
			IdleTimeout:        7 * time.Second,
			IdleCheckFrequency: 5 * time.Second,
		}),
		keepAliveInterval: 2 * time.Second,
	}
	go func() {
		redisClient.doKeepAliveInterval(panicIfDisconnect)
	}()

	lastRdsClient = redisClient
	return lastRdsClient
}

func (p *RedisClient) doKeepAliveInterval(panicIfDisconnect bool) {
	for {
		<-time.After(p.keepAliveInterval)
		statCmd := p.session.Ping()
		if nil != statCmd && nil != statCmd.Err() && panicIfDisconnect {
			fmt.Printf("redis keep alived failed:%v\n", statCmd.Err())
			panic(statCmd.Err)
		}
	}
}

/**
 * 模糊查询 key
 */
func (p *RedisClient) Keys(pattern string) (vals []string, err error) {
	result := p.session.Keys(pattern)
	err = result.Err()
	vals = result.Val()
	return
}

/**
 * 批量设置
 * vals : len(vals) 必定等于 len(keys)
 * err : redis.Nil, key对应的数据不存在
 */
func (p *RedisClient) MSet(vals []interface{}) (string, error) {
	result := p.session.MSet(vals...)
	return result.Result()
}

/**
 * 批量获取
 * vals : len(vals) 必定等于 len(keys)
 * err : redis.Nil, key对应的数据不存在
 */
func (p *RedisClient) Get(keys []string) (vals []string, err error) {
	// *StringCmd
	vals = make([]string, len(keys))
	for idx := range keys {
		stringCmd := p.session.Get(keys[idx])
		if nil != stringCmd.Err() {
			err = stringCmd.Err()
		}
		vals[idx] = stringCmd.Val()
	}
	return
}

func (p *RedisClient) MGet(keys []string) (vals []interface{}, err error) {
	Cmd := p.session.MGet(keys...)
	return Cmd.Result()
}

/**
 * 设置
 * expiration : Zero expiration means the key has no expiration time.
 *
 */
func (p *RedisClient) Set(key string, val string, second int) (result string, err error) {
	intCmd := p.session.Set(key, val, time.Duration(second)*time.Second)
	result, err = intCmd.Result()
	return
}

func (p *RedisClient) SetAdd(key string, vals []interface{}) (result int64, err error) {
	intCmd := p.session.SAdd(key, vals...)
	result, err = intCmd.Result()
	return
}

func (p *RedisClient) SetNx(key string, val string, second int) (result bool, err error) {
	intCmd := p.session.SetNX(key, val, time.Duration(second)*time.Second)
	result, err = intCmd.Result()
	return
}

func (c *RedisClient) HSet(key string, field string, val interface{}) (bool, error) {
	return c.session.HSet(key, field, val).Result()
}

func (c *RedisClient) HDel(key string, fields ...string) (int64, error) {
	return c.session.HDel(key, fields...).Result()
}

func (c *RedisClient) HGet(key string, field string) (string, error) {
	return c.session.HGet(key, field).Result()
}

func (c *RedisClient) HGetAll(key string) (map[string]string, error) {
	return c.session.HGetAll(key).Result()
}

func (c *RedisClient) HKeys(key string) ([]string, error) {
	return c.session.HKeys(key).Result()
}

/**
 * 单个获取
 * err : { key对应的数据不存在 => redis.Nil, 其他redis错误 => err }
 */
func (p *RedisClient) GetOne(key string) (string, error) {
	stringCmd := p.session.Get(key)
	result, err := stringCmd.Result()
	return result, err
}

/**
 * 单个删除
 * int64 : >0, 成功删除1个Key对应的数据
 * err : if nil ==> 成功提交redis并成功执行返回
 */
func (p *RedisClient) DelOne(key string) (result int64, err error) {
	intCmd := p.session.Del(key)
	result, err = intCmd.Result()
	return
}

/**
 * 队列压入
 */
func (p *RedisClient) LPush(key string, val string) (result int64, err error) {
	intCmd := p.session.LPush(key, val)
	result, err = intCmd.Result()
	return
}

/**
 * 队列抛出
 */
func (p *RedisClient) RPop(key string) (result string, err error) {
	intCmd := p.session.RPop(key)
	result, err = intCmd.Result()
	return
}

/**
 * 队列长度
 */
func (p *RedisClient) LLen(key string) (result int64, err error) {
	intCmd := p.session.LLen(key)
	result, err = intCmd.Result()
	return
}

/**
 * 队列裁剪
 */
func (p *RedisClient) LTrim(key string, start int64, end int64) (result string, err error) {
	intCmd := p.session.LTrim(key, start, end)
	result, err = intCmd.Result()
	return
}

/**
 * 获取队列里多个值
 */
func (p *RedisClient) LRange(key string, start int64, end int64) (result []string, err error) {
	stringsliceCmd := p.session.LRange(key, start, end)
	result, err = stringsliceCmd.Result()
	return
}

/**
 * bool : true, key对应的数据得到了更新
 *        false，没有key对应的数据
 * err : if nil ==> 成功提交redis并成功执行返回
 */
func (p *RedisClient) Expire(key string, expiration time.Duration) (result bool, err error) {
	boolCmd := p.session.Expire(key, expiration)
	result, err = boolCmd.Result()
	return
}

/**
 * int64 : >0, key对应的数值型数据incr后的结果
 *          0, key对应的数据非数值
 * err : if nil ==> 成功提交redis并成功执行返回
 */
func (p *RedisClient) Publish(ch string, msg string) (result int64, err error) {
	intCmd := p.session.Publish(ch, msg)
	result, err = intCmd.Result()
	return
}

/**
 * int64 : >0, key对应的数值型数据incr后的结果
 *          0, key对应的数据非数值
 * err : if nil ==> 成功提交redis并成功执行返回 or key对应的数据非数值
 */
func (p *RedisClient) Incr(key string) (result int64, err error) {
	intCmd := p.session.Incr(key)
	result, err = intCmd.Result()
	return
}

func (p *RedisClient) TTL(key string) (result time.Duration, err error) {
	intCmd := p.session.TTL(key)
	result, err = intCmd.Result()
	return
}

func (p *RedisClient) SRem(key string, members []interface{}) (result int64, err error) {
	intCmd := p.session.SRem(key, members...)
	result, err = intCmd.Result()
	return
}

func (p *RedisClient) SMembers(key string) (result []string, err error) {
	stringCmd := p.session.SMembers(key)
	result, err = stringCmd.Result()
	return
}
