package db

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	lastRdsClient *RedisClient
)

// RedisClient ...
type RedisClient struct {
	addr              string
	port              string
	passwd            string
	db                int
	session           *redis.Client
	ctx               context.Context
	keepAliveInterval time.Duration //  second
	timeout           time.Duration
}

// IsNotExistError ...
func IsNotExistError(err error) bool {
	if err == redis.Nil {
		return true
	}
	return false
}

// GetLastRedisClient ...
func GetLastRedisClient() *RedisClient {
	return lastRdsClient
}

// NewRedisClient ...
func NewRedisClient(addr string, passwd string, db int, panicIfDisconnect bool) *RedisClient {
	redisClient := &RedisClient{
		addr:   addr,
		passwd: passwd,
		db:     db,
		session: redis.NewClient(&redis.Options{
			Addr:               addr,
			Password:           passwd,
			DB:                 db,
			PoolSize:           0xFF,
			IdleTimeout:        7 * time.Second,
			IdleCheckFrequency: 5 * time.Second,
		}),
		keepAliveInterval: 2 * time.Second,
		ctx:               context.Background(),
		timeout:           2 * time.Second,
	}
	go func() {
		redisClient.doKeepAliveInterval(panicIfDisconnect)
	}()

	lastRdsClient = redisClient
	return lastRdsClient
}

func (r *RedisClient) doKeepAliveInterval(panicIfDisconnect bool) {
	for {
		<-time.After(r.keepAliveInterval)
		ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
		defer cancel()
		statCmd := r.session.Ping(ctx)
		if nil != statCmd && nil != statCmd.Err() && panicIfDisconnect {
			fmt.Printf("redis keep alived failed:%v\n", statCmd.Err())
			panic(statCmd.Err)
		}
	}
}

// Keys FuzzyQuery
func (r *RedisClient) Keys(pattern string) (vals []string, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	result := r.session.Keys(ctx, pattern)
	err = result.Err()
	vals = result.Val()
	return
}

// MSet Batch Set
func (r *RedisClient) MSet(vals []interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	result := r.session.MSet(ctx, vals...)
	return result.Result()
}

// Get Batch Get
func (r *RedisClient) Get(key string) (val string, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	cmd := r.session.Get(ctx, key)
	if nil != cmd.Err() {
		err = cmd.Err()
	}
	val = cmd.Val()
	return
}

// MGet ...
func (r *RedisClient) MGet(keys []string) (vals []interface{}, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	Cmd := r.session.MGet(ctx, keys...)
	return Cmd.Result()
}

// Set ...
func (r *RedisClient) Set(key string, val string, expire int) (result string, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	intCmd := r.session.Set(ctx, key, val, time.Duration(expire)*time.Second)
	result, err = intCmd.Result()
	return
}

// SetAdd ...
func (r *RedisClient) SetAdd(key string, vals []interface{}) (result int64, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	intCmd := r.session.SAdd(ctx, key, vals...)
	result, err = intCmd.Result()
	return
}

// SetNx ...
func (r *RedisClient) SetNx(key string, val string, expire int) (result bool, err error) {
	// result mean whether the key exist.
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	intCmd := r.session.SetNX(ctx, key, val, time.Duration(expire)*time.Second)
	result, err = intCmd.Result()
	return
}

// HSet ...
func (r *RedisClient) HSet(key string, field string, val interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	return r.session.HSet(ctx, key, field, val).Result()
}

// HDel ...
func (r *RedisClient) HDel(key string, fields ...string) (int64, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	return r.session.HDel(ctx, key, fields...).Result()
}

// HGet ...
func (r *RedisClient) HGet(key string, field string) (string, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	return r.session.HGet(ctx, key, field).Result()
}

// HGetAll ...
func (r *RedisClient) HGetAll(key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	return r.session.HGetAll(ctx, key).Result()
}

// HKeys ...
func (r *RedisClient) HKeys(key string) ([]string, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	return r.session.HKeys(ctx, key).Result()
}

// GetOne ...
func (r *RedisClient) GetOne(key string) (string, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	stringCmd := r.session.Get(ctx, key)
	result, err := stringCmd.Result()
	return result, err
}

// DelOne ...
func (r *RedisClient) DelOne(key string) (result int64, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	intCmd := r.session.Del(ctx, key)
	result, err = intCmd.Result()
	return
}

// LPush ...
func (r *RedisClient) LPush(key string, val string) (result int64, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	intCmd := r.session.LPush(ctx, key, val)
	result, err = intCmd.Result()
	return
}

// RPop ...
func (r *RedisClient) RPop(key string) (result string, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	intCmd := r.session.RPop(ctx, key)
	result, err = intCmd.Result()
	return
}

// LLen ...
func (r *RedisClient) LLen(key string) (result int64, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	intCmd := r.session.LLen(ctx, key)
	result, err = intCmd.Result()
	return
}

// LTrim ...
func (r *RedisClient) LTrim(key string, start int64, end int64) (result string, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	intCmd := r.session.LTrim(ctx, key, start, end)
	result, err = intCmd.Result()
	return
}

// LRange ...
func (r *RedisClient) LRange(key string, start int64, end int64) (result []string, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	stringsliceCmd := r.session.LRange(ctx, key, start, end)
	result, err = stringsliceCmd.Result()
	return
}

// Expire ...
func (r *RedisClient) Expire(key string, expiration time.Duration) (result bool, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	boolCmd := r.session.Expire(ctx, key, expiration)
	result, err = boolCmd.Result()
	return
}

// Publish ...
func (r *RedisClient) Publish(ch string, msg string) (result int64, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	intCmd := r.session.Publish(ctx, ch, msg)
	result, err = intCmd.Result()
	return
}

// Incr ...
func (r *RedisClient) Incr(key string) (result int64, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	intCmd := r.session.Incr(ctx, key)
	result, err = intCmd.Result()
	return
}

// TTL ...
func (r *RedisClient) TTL(key string) (result time.Duration, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	intCmd := r.session.TTL(ctx, key)
	result, err = intCmd.Result()
	return
}

// SRem ...
func (r *RedisClient) SRem(key string, members []interface{}) (result int64, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	intCmd := r.session.SRem(ctx, key, members...)
	result, err = intCmd.Result()
	return
}

// SMembers ...
func (r *RedisClient) SMembers(key string) (result []string, err error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	stringCmd := r.session.SMembers(ctx, key)
	result, err = stringCmd.Result()
	return
}
