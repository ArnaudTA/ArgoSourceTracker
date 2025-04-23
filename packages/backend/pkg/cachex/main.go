package cachex

import (
	"argocd-watcher/pkg/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

var ttl = time.Duration(config.Global.RegistryCacheTTL) * time.Minute

type CacheService[T any] struct {
	ctx        context.Context
	client     *cache.Cache
	redis      *redis.Client
	localCache cache.LocalCache
	channel    string
	useRedis   bool
}

func NewCacheService[T any](ctx context.Context, channel string) *CacheService[T] {
	useRedis := config.Global.Redis.Enabled
	addr := fmt.Sprintf("%s:%d", config.Global.Redis.Host, config.Global.Redis.Port)

	var redisClient *redis.Client
	if useRedis {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: config.Global.Redis.Password,
			DB:       config.Global.Redis.Db,
		})
	}

	var redisBackend redis.UniversalClient
	if useRedis {
		redisBackend = redisClient
	}
	local := cache.NewTinyLFU(config.Global.RegistryCacheTTL, time.Minute)
	c := cache.New(&cache.Options{
		Redis:      redisBackend,
		LocalCache: local,
	})

	svc := &CacheService[T]{
		ctx:        ctx,
		client:     c,
		redis:      redisClient,
		channel:    channel,
		localCache: local,
		useRedis:   useRedis,
	}

	if useRedis {
		go svc.listenInvalidations()
	}

	return svc
}

func (c *CacheService[T]) Set(key string, value T) error {
	err := c.client.Set(&cache.Item{
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
	if err == nil && c.useRedis {
		// Publie l'invalidation (simple payload = key)
		_ = c.redis.Publish(c.ctx, c.channel, key)
	}
	return err
}

func (c *CacheService[T]) Get(key string) (T, error) {
	var val T
	err := c.client.Get(c.ctx, key, &val)
	return val, err
}

func (c *CacheService[T]) Delete(key string) error {
	err := c.client.Delete(c.ctx, key)
	if err == nil && c.useRedis {
		_ = c.redis.Publish(c.ctx, c.channel, key)
	}
	return err
}

func (c *CacheService[T]) listenInvalidations() {
	sub := c.redis.Subscribe(c.ctx, c.channel)
	ch := sub.Channel()

	for msg := range ch {
		key := msg.Payload
		c.localCache.Del(key) // ✅ seulement local
		log.Printf("[CacheService] Invalidation received: %s\n", key)
	}
}
