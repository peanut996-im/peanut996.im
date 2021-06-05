// Package encoding
// @Title  encoding.go
// @Description  提供需要的简化的编码操作
// @Author  peanut996
// @Update  peanut996  2021/5/22 1:47
package tool

import (
	"crypto/sha1"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"time"
)

func EncryptBySha1(plain string) string {
	h := sha1.New()
	h.Write([]byte(plain))
	return fmt.Sprintf("%X", h.Sum(nil))
}

func MapToStruct(input, output interface{}) error {
	return mapstructure.Decode(input, output)
}

func ToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
		// Convert it by parsing
	}
}

func Decode(input map[string]interface{}, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			ToTimeHookFunc()),
		Result: result,
	})
	if err != nil {
		return err
	}

	if err := decoder.Decode(input); err != nil {
		return err
	}
	return err
}
