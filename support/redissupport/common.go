package redissupport

import (
	"context"
	"crawler-distributed/config"
	"github.com/redis/go-redis/v9"
	"log"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func SaveToRedis(client *redis.Client, key string, value string, ctx context.Context) {
	err := client.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Fatalf("[SaveToRedis] error save key to redissupport:%v", err)
	}
}

func GetFromRedis(client *redis.Client, key string, ctx context.Context) string {
	// 当待查询的 key 不存在时，redis 会返回 redis:nil 的错误信息
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return result
}
