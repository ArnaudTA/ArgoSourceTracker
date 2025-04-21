package cache

import (
	"argocd-watcher/pkg/config"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var ctx = context.Background()

var rdbClient *redis.Client

const invalidationChannel = "cache-invalidation"

var (
	memoryCache = make(map[string][]byte)
	mu          sync.RWMutex
	pubsubOnce  sync.Once
)

func GetClient() *redis.Client {
	if rdbClient == nil {
		addr := fmt.Sprintf("%s:%d", config.Global.Redis.Host, config.Global.Redis.Port)
		rdbClient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: config.Global.Redis.Password,
			DB:       config.Global.Redis.Db,
		})
	}
	return rdbClient
}

func Store(key string, value []byte, ttl time.Duration) {
	if ttl == 0 {
		ttl = 5 * time.Minute
	}
	// En mémoire
	mu.Lock()
	memoryCache[key] = value
	mu.Unlock()

	client := GetClient()
	client.Set(ctx, key, value, ttl)
}

// Load tente d'abord la mémoire, puis Redis.
func Load(key string) ([]byte, error) {
	// In-memory first
	mu.RLock()
	val, ok := memoryCache[key]
	mu.RUnlock()
	if ok {
		logrus.Tracef("Cache retieved from memory: %s", key)
		return val, nil
	}

	// Fallback Redis
	client := GetClient()
	data, err := client.Get(ctx, key).Bytes()
	if err != nil {
		logrus.Tracef("Cache not retieved: %s", key)
		return nil, err
	}

	// Synchronise dans la mémoire
	mu.Lock()
	memoryCache[key] = data
	mu.Unlock()
	
	logrus.Tracef("Cache retieved from redis: %s", key)
	return data, nil
}

func Delete(key string) {
	mu.Lock()
	delete(memoryCache, key)
	mu.Unlock()

	client := GetClient()
	client.Del(ctx, key)

	broadcastInvalidation(key)
}

func broadcastInvalidation(key string) {
	client := GetClient()
	client.Publish(ctx, invalidationChannel, key)
}

func startInvalidationListener() {
	pubsubOnce.Do(func() {
		go func() {
			client := GetClient()
			pubsub := client.Subscribe(ctx, invalidationChannel)
			ch := pubsub.Channel()
			for msg := range ch {
				mu.Lock()
				delete(memoryCache, msg.Payload)
				mu.Unlock()
			}
		}()
	})
}

func Init() {
	startInvalidationListener()
}