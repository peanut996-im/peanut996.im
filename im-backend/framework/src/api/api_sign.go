package api

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"sort"
)

func CheckSignFromJsonParams(s string, appKey string) (bool, error) {

	j := make(map[string]interface{})

	err := json.Unmarshal([]byte(s), &j)
	getSign, ok := j["sign"]
	if !ok || getSign == "" || err != nil {
		return false, err
	}

	makeSign, err := MakeSignWithJsonParams(s, appKey)
	if err != nil {
		return false, err
	}

	return makeSign == getSign, nil
}

func CheckSignFromQueryParams(params url.Values, appKey string) (bool, error) {
	getSign := params.Get("sign")
	if getSign == "" {
		return false, nil
	}

	makeSign, err := MakeSignWithQueryParams(params, appKey)
	if err != nil {
		return false, err
	}

	return makeSign == getSign, nil
}

func MakeSignWithQueryParams(params url.Values, appKey string) (string, error) {
	// To store the keys in slice in sorted order
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha1.New()
	fmt.Print("sign before sha1: ")
	for _, k := range keys {
		switch k {
		case "sign":
			continue
		case "EIO":
			continue
		case "transport":
			continue
		case "t":
			continue
		}
		fmt.Printf("%v%v", k, params.Get(k))
		if _, err := io.WriteString(h, fmt.Sprintf("%v%v", k, params.Get(k))); err != nil {
			return "", err
		}
	}
	fmt.Println(appKey)
	if _, err := io.WriteString(h, appKey); err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", h.Sum(nil)), nil

}

//func MakeSignWithJsonParams(object interface{}, appkey string) (string, error) {
//	getType := reflect.TypeOf(object)
//	getValue := reflect.ValueOf(object)
//	if getType.Kind() == reflect.Ptr {
//		getType = getType.Elem()
//		getValue = getValue.Elem()
//	}
//
//	vals := url.Values{}
//	for i := 0; i < getType.NumField(); i++ {
//		field := getType.Field(i)
//		value := getValue.Field(i)
//		tag := field.Tag.Get("json")
//		if strings.Compare("sign", tag) == 0 {
//			continue
//		}
//
//		switch value.Kind() {
//		case reflect.Ptr:
//		case reflect.Struct:
//		case reflect.Array:
//		case reflect.Map:
//		case reflect.UnsafePointer:
//		case reflect.Slice:
//		default:
//			vals.Add(tag, fmt.Sprintf("%v", value))
//		}
//	}
//
//	return MakeSignWithQueryParams(vals, appkey)
//}

func MakeSignWithJsonParams(s string, appkey string) (string, error) {
	j := make(map[string]interface{})
	vals := url.Values{}

	err := json.Unmarshal([]byte(s), &j)

	if err != nil {
		return "", err
	}

	for k, v := range j {
		switch k {
		case "sign":
			continue
		case "EIO":
			continue
		case "transport":
			continue
		case "t":
			continue
		}

		value := reflect.ValueOf(v)
		switch value.Kind() {
		case reflect.Ptr:
		case reflect.Struct:
		case reflect.Array:
		case reflect.Map:
		case reflect.UnsafePointer:
		case reflect.Slice:
		default:
			vals.Add(k, fmt.Sprintf("%v", value))
		}

	}
	return MakeSignWithQueryParams(vals, appkey)
}
