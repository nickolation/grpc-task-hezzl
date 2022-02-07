package cache

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/nickolation/grpc-task-hezzl/configs"
	"github.com/nickolation/grpc-task-hezzl/model"
	"github.com/sirupsen/logrus"
)

type GrpcCacheStorage interface {
	SetToCasheUserList(ctx context.Context, list []model.User) error
	GetCachedUserList(ctx context.Context) (string, error)
}

type CacheStorage struct {
	Storage *cache.Cache
	Config  *configs.GrpcGonfigManager
}

const _defaultRedisPassword = ""
const _defaultRedisDbId = 0

func (cs *CacheStorage) connectToRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cs.Config.GetRedisAddr(),
		Password: _defaultRedisPassword,
		DB:       _defaultRedisDbId,
	})
}

func (cs *CacheStorage) defaultConnectToRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cs.Config.GetRedisAddr(),
		Password: _defaultRedisPassword,
		DB:       _defaultRedisDbId,
	})
}

func NewGrpcCacheStorage(cfg *configs.GrpcGonfigManager) *CacheStorage {
	cs := new(CacheStorage)
	cs.Config = cfg
	return &CacheStorage{
		Storage: cache.New(&cache.Options{
			Redis: cs.connectToRedis(),
		}),
	}
}

const _lastCasheIdKey = "last"

var _defaultCacheTTL = time.Minute

func (cs *CacheStorage) SetToCasheUserList(ctx context.Context, list []model.User) error {
	return cs.Storage.Set(&cache.Item{
		Ctx: ctx,
		Key: _lastCasheIdKey,
		Value: func() []byte {
			b, err := marshalUserList(list)
			if err != nil {
				logrus.Errorf("marshaling userList err - [%e]", err)
			}

			return b
		}(),
		TTL: _defaultCacheTTL,
	})
}

func (cs *CacheStorage) GetCachedUserList(ctx context.Context) (string, error) {
	var s string
	err :=  cs.Storage.Get(ctx, _lastCasheIdKey, &s)
	return s, err
}

