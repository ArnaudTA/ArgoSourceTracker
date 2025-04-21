package cache

import (
	"argocd-watcher/pkg/config"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var rdbClient *redis.Client

func GetClient() *redis.Client {
	if rdbClient == nil {
		addr := fmt.Sprintf("%s:%d", config.Global.Redis.Host, config.Global.Redis.Port)
		rdbClient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: config.Global.Redis.Password,
			DB:       10, // use default DB
		})
	}
	return rdbClient
}
func Store(key string, value []byte, ttl time.Duration) {
	if ttl == 0 {
		ttl = 5 * time.Minute
	}
	client := GetClient()

	client.Set(ctx, key, value, ttl)
}

func Load(key string) ([]byte, error) {
	client := GetClient()

	return client.Get(ctx, key).Bytes()

}

func Delete(key string) {
	client := GetClient()

	client.Del(ctx, key)
}
