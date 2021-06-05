// Package api
// @Title  api_token.go
// @Description  提供对token的常用操作
// @Author  peanut996
// @Update  peanut996  2021/5/22 1:43
package api

import (
	"fmt"
	"framework/api/model"
	"framework/tool"

	"framework/db"
	"time"
)

const (
	//UidToTokenFormat redis key name for value: token
	UidToTokenFormat = "%v_to_token"
	//TokenToUidFormat redis key name for value: uid
	TokenToUidFormat = "%v_to_uid"
	//DefaultTokenExpireTime Token default expiration time
	DefaultTokenExpireTime = 6 * 60 * 60
)

func CheckToken(token string) (*model.User, error) {
	rds := db.GetLastRedisClient()

	uid, err := rds.Get(TokenToUIDFormat(token))

	if err != nil {
		return nil, err
	}

	// 重置时间
	go ResetTokenTime(token,uid)

	return model.GetUserByUID(uid)
}

//InsertToken token插入数据库，若已存在则直接返回已存在的token
func InsertToken(uid string) (string, error) {
	// 先查询 已存在则直接获取数据库token 并重新设置过期时间
	rds := db.GetLastRedisClient()
	tokenKey := UIDToTokenFormat(uid)

	token, err := rds.Get(tokenKey)
	if err != nil {
		if db.IsNotExistError(err) {
			token = GenerateToken(uid)
		} else {
			// redis 有问题
			return "", err
		}
	}

	// 已存在 重置
	// redis uid => token
	go ResetTokenTime(token,uid)
	return token, nil
}

//GenerateToken 根据uid和时间戳生成token
func GenerateToken(uid string) string {
	ts := time.Now().Unix()
	origin := fmt.Sprintf("%v_%v", uid, ts)
	return tool.EncryptBySha1(origin)
}

func DeleteToken(token string) error {
	rds := db.GetLastRedisClient()
	// 只需要删除token => uid 即可使token失效 uid=>token可复用
	_, err := rds.DelOne(TokenToUIDFormat(token))

	if nil != err {
		return err
	}
	return nil
}

func ResetTokenTime(token,uid string) error {
	rds := db.GetLastRedisClient()
	tokenKey := UIDToTokenFormat(uid)
	uidKey := TokenToUIDFormat(token)
	_, err := rds.Set(tokenKey, token, DefaultTokenExpireTime)
	if nil != err {
		return  err
	}
	// redis token => uid
	_, err = rds.Set(uidKey, uid, DefaultTokenExpireTime)
	if nil != err {
		return  err
	}
	return nil
}

func UIDToTokenFormat(uid string) string {
	return fmt.Sprintf(UidToTokenFormat, uid)
}

func TokenToUIDFormat(token string) string {
	return fmt.Sprintf(TokenToUidFormat, token)
}
