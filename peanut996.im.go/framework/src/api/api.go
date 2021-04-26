package api

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"sort"
	"strings"
)

func MakeSign(params url.Values, appKey string) (string, error) {
	// To store the keys in slice in sorted order
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha1.New()
	for _, k := range keys {
		switch k {
		case "sign":
			continue
		case "EIO":
			continue
		case "transport":
			continue
		}
		fmt.Printf("%v%v", k, params.Get(k))
		if _, err := io.WriteString(h, fmt.Sprintf("%v%v", k, params.Get(k))); err != nil {
			return "", err
		}
	}
	fmt.Print(appKey)
	if _, err := io.WriteString(h, appKey); err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", h.Sum(nil)), nil

}

func MakeSignWithJsonParams(object interface{}, appkey string) (string, error) {
	getType := reflect.TypeOf(object)
	getValue := reflect.ValueOf(object)
	if getType.Kind() == reflect.Ptr {
		getType = getType.Elem()
		getValue = getValue.Elem()
	}

	vals := url.Values{}
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i)
		tag := field.Tag.Get("json")
		if strings.Compare("sign", tag) == 0 {
			continue
		}

		switch value.Kind() {
		case reflect.Ptr:
		case reflect.Struct:
		case reflect.Array:
		case reflect.Map:
		case reflect.UnsafePointer:
		case reflect.Slice:
		default:
			vals.Add(tag, fmt.Sprintf("%v", value))
		}
	}

	return MakeSign(vals, appkey)
}
