package db

import (
	"framework/cfgargs"
	"reflect"
	"testing"
	"time"
)

var (
	redisClient *RedisClient
	redisConfig *cfgargs.SrvConfig
)

func init() {
	cfg, err := cfgargs.GetSrvConfig("../../etc/config-example.yaml")
	if nil != err {
		panic("get config error")
	}
	client := NewRedisClient(cfgargs.GetRedisAddr(cfg), cfg.Redis.Passwd, cfg.Redis.DB, true)
	redisClient = client
	redisConfig = cfg
}

func TestGetLastRedisClient(t *testing.T) {
	tests := []struct {
		name string
		want *RedisClient
	}{
		{"case0", redisClient},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLastRedisClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLastRedisClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRedisClient(t *testing.T) {
	type args struct {
		addr              string
		pass              string
		db                int
		panicIfDisconnect bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"case0", args{redisConfig.Redis.Host, redisConfig.Redis.Passwd, redisConfig.Redis.DB, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRedisClient(tt.args.addr, tt.args.pass, tt.args.db, tt.args.panicIfDisconnect); got == nil {
				t.Errorf("NewRedisClient() = nil")
			}
		})
	}
}

func TestRedisClient_Keys(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name     string
		r        *RedisClient
		args     args
		wantVals []string
		wantErr  bool
	}{
		{"case0", redisClient, args{"user*"}, []string{"user0", "user1", "user2"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVals, err := tt.r.Keys(tt.args.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.Keys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVals, tt.wantVals) {
				t.Errorf("RedisClient.Keys() = %v, want %v", gotVals, tt.wantVals)
			}
		})
	}
}

func TestRedisClient_MSet(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	vals := []interface{}{
		"user0", 545,
		"user1", 32,
		"user2", 235,
	}
	tests := []struct {
		name    string
		r       *RedisClient
		args    args
		want    string
		wantErr bool
	}{
		{"case0", redisClient, args{vals}, "OK", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.MSet(tt.args.vals)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.MSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisClient.MSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisClient_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		r       *RedisClient
		args    args
		wantVal string
		wantErr bool
	}{
		{"case0", redisClient, args{"user0"}, "545", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVals, err := tt.r.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVals, tt.wantVal) {
				t.Errorf("RedisClient.Get() = %v, want %v", gotVals, tt.wantVal)
			}
		})
	}
}

func TestRedisClient_MGet(t *testing.T) {
	type args struct {
		keys []string
	}
	tests := []struct {
		name     string
		r        *RedisClient
		args     args
		wantVals []interface{}
		wantErr  bool
	}{
		{"case0", redisClient, args{[]string{"user0", "user1", "user2"}}, []interface{}{"545", "32", "235"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVals, err := tt.r.MGet(tt.args.keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.MGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVals, tt.wantVals) {
				t.Errorf("RedisClient.MGet() = %v, want %v", gotVals, tt.wantVals)
			}
		})
	}
}

func TestRedisClient_Set(t *testing.T) {
	type args struct {
		key    string
		val    string
		expire int
	}
	tests := []struct {
		name       string
		r          *RedisClient
		args       args
		wantResult string
		wantErr    bool
	}{
		{"case0", redisClient, args{"key", "val", 0}, "OK", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.Set(tt.args.key, tt.args.val, tt.args.expire)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("RedisClient.Set() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestRedisClient_SetNx(t *testing.T) {
	type args struct {
		key    string
		val    string
		expire int
	}
	tests := []struct {
		name       string
		r          *RedisClient
		args       args
		wantResult bool
		wantErr    bool
	}{
		{"case0", redisClient, args{"key", "val", 0}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.SetNx(tt.args.key, tt.args.val, tt.args.expire)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.SetNx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("RedisClient.SetNx() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestRedisClient_HSet(t *testing.T) {
	type args struct {
		key   string
		field string
		val   interface{}
	}
	tests := []struct {
		name    string
		r       *RedisClient
		args    args
		want    int64
		wantErr bool
	}{
		{"case0", redisClient, args{"hmap", "key", "val"}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.HSet(tt.args.key, tt.args.field, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.HSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisClient.HSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisClient_HGet(t *testing.T) {
	type args struct {
		key   string
		field string
	}
	tests := []struct {
		name    string
		r       *RedisClient
		args    args
		want    string
		wantErr bool
	}{
		{"case0", redisClient, args{"hmap", "key"}, "val", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.HGet(tt.args.key, tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.HGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisClient.HGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisClient_HGetAll(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		r       *RedisClient
		args    args
		want    map[string]string
		wantErr bool
	}{
		{"case0", redisClient, args{"hmap"}, map[string]string{"key": "val"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.HGetAll(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.HGetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisClient.HGetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisClient_HKeys(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		r       *RedisClient
		args    args
		want    []string
		wantErr bool
	}{
		{"case0", redisClient, args{"hmap"}, []string{"key"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.HKeys(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.HKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisClient.HKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestRedisClient_HDel(t *testing.T) {
	type args struct {
		key    string
		fields []string
	}
	tests := []struct {
		name    string
		r       *RedisClient
		args    args
		want    int64
		wantErr bool
	}{
		{"case0", redisClient, args{"hmap", []string{"key"}}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.HDel(tt.args.key, tt.args.fields...)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.HDel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisClient.HDel() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestRedisClient_GetOne(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		r       *RedisClient
		args    args
		want    string
		wantErr bool
	}{
		{"case0", redisClient, args{"user0"}, "545", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetOne(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.GetOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisClient.GetOne() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisClient_DelOne(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name       string
		r          *RedisClient
		args       args
		wantResult int64
		wantErr    bool
	}{
		{"case0", redisClient, args{"key"}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.DelOne(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.DelOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("RedisClient.DelOne() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestRedisClient_LPush(t *testing.T) {
	type args struct {
		key string
		val string
	}
	tests := []struct {
		name       string
		r          *RedisClient
		args       args
		wantResult int64
		wantErr    bool
	}{
		{"case0", redisClient, args{"list", "val1"}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.LPush(tt.args.key, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.LPush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("RedisClient.LPush() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestRedisClient_LLen(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name       string
		r          *RedisClient
		args       args
		wantResult int64
		wantErr    bool
	}{
		{"case0", redisClient, args{"list"}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.LLen(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.LLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("RedisClient.LLen() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestRedisClient_LRange(t *testing.T) {
	type args struct {
		key   string
		start int64
		end   int64
	}
	tests := []struct {
		name       string
		r          *RedisClient
		args       args
		wantResult []string
		wantErr    bool
	}{
		{"case0", redisClient, args{"list", 0, 1}, []string{"val1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.LRange(tt.args.key, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.LRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("RedisClient.LRange() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestRedisClient_LTrim(t *testing.T) {
	type args struct {
		key   string
		start int64
		end   int64
	}
	tests := []struct {
		name       string
		r          *RedisClient
		args       args
		wantResult string
		wantErr    bool
	}{
		{"case0", redisClient, args{"list", 0, 1}, "OK", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.LTrim(tt.args.key, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.LTrim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("RedisClient.LTrim() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestRedisClient_RPop(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name       string
		r          *RedisClient
		args       args
		wantResult string
		wantErr    bool
	}{
		{"case0", redisClient, args{"list"}, "val1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.RPop(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.RPop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("RedisClient.RPop() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestRedisClient_TTL(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name       string
		r          *RedisClient
		args       args
		wantResult time.Duration
		wantErr    bool
	}{
		{"case0", redisClient, args{"key"}, -2 * time.Nanosecond, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.TTL(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.TTL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("RedisClient.TTL() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestRedisClient_SetAdd(t *testing.T) {
	type args struct {
		key  string
		vals []interface{}
	}
	tests := []struct {
		name       string
		r          *RedisClient
		args       args
		wantResult int64
		wantErr    bool
	}{
		{"case0", redisClient, args{"SetKey", []interface{}{1, 2, 3, 4}}, 4, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.SetAdd(tt.args.key, tt.args.vals)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.SetAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("RedisClient.SetAdd() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestRedisClient_SMembers(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name       string
		r          *RedisClient
		args       args
		wantResult []string
		wantErr    bool
	}{
		{"case0", redisClient, args{"SetKey"}, []string{"1", "2", "3", "4"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.SMembers(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.SMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("RedisClient.SMembers() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestRedisClient_SRem(t *testing.T) {
	type args struct {
		key     string
		members []interface{}
	}
	tests := []struct {
		name       string
		r          *RedisClient
		args       args
		wantResult int64
		wantErr    bool
	}{
		{"case0", redisClient, args{"SetKey", []interface{}{1, 2, 3, 4}}, 4, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.SRem(tt.args.key, tt.args.members)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.SRem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("RedisClient.SRem() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestRedisClient_Expire(t *testing.T) {
	type args struct {
		key        string
		expiration time.Duration
	}
	tests := []struct {
		name       string
		r          *RedisClient
		args       args
		wantResult bool
		wantErr    bool
	}{
		{"case0", redisClient, args{"hmap", 1 * time.Second}, false, false},
		{"case1", redisClient, args{"SetKey", 1 * time.Second}, false, false},
		{"case2", redisClient, args{"list", 1 * time.Second}, false, false},
		{"case2", redisClient, args{"key", 1 * time.Second}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.Expire(tt.args.key, tt.args.expiration)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisClient.Expire() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("RedisClient.Expire() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
