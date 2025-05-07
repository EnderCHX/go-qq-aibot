package df

import (
	"context"
	"fmt"
	"github.com/EnderCHX/go-qq-aibot/config"
	"github.com/redis/go-redis/v9"
	"time"
)

var ctx = context.Background()
var rdb *redis.Client

func init() {
	cf := config.GetConfig()
	rdb = redis.NewClient(&redis.Options{
		Addr:     cf.Redis.Host + ":" + cf.Redis.Port,
		Password: cf.Redis.Password,
		DB:       cf.Redis.DB,
	})
}

func reConnectRedis() bool {
	if rdb == nil {
		timeout := time.Second * 10
		timenow := time.Now()
		for rdb == nil && time.Since(timenow) < timeout {
			rdb = redis.NewClient(&redis.Options{
				Addr:     config.GetConfig().Redis.Host + ":" + config.GetConfig().Redis.Port,
				Password: config.GetConfig().Redis.Password,
				DB:       config.GetConfig().Redis.DB,
			})
		}
	}
	return rdb != nil
}
func SaveCookie(cookie, qqid string) error {
	if !reConnectRedis() {
		return fmt.Errorf("connect to redis timeout")
	}
	return rdb.Set(ctx, "df:cookie:"+qqid, cookie, 0).Err()
}

func GetCookie(qqid string) (string, error) {
	if !reConnectRedis() {
		return "", fmt.Errorf("connect to redis timeout")
	}
	return rdb.Get(ctx, "df:cookie:"+qqid).Result()
}
