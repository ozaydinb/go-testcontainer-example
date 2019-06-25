package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/ozaydinb/go-testcontainer-example/model"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RedisCacheManager interface {
	GetString(key string) (string, error)
	SetString(key string, value interface{}) error
}

type RedisCacheManagerImp struct {
	pool        *redis.Pool
	redisConfig model.RedisConfig
}

func NewRedisCacheManager(redisConfig model.RedisConfig) RedisCacheManager {
	redisPool := initRedisPool(redisConfig)
	cleanupHook(redisPool)
	return &RedisCacheManagerImp{
		pool:        redisPool,
		redisConfig: redisConfig,
	}
}

// private functions

func initRedisPool(redisConfig model.RedisConfig) *redis.Pool {
	redisHost := fmt.Sprintf("%s:%v", redisConfig.Host, redisConfig.Port)
	return &redis.Pool{

		MaxIdle:     redisConfig.MaxIdle,
		IdleTimeout: redisConfig.IdleTimeoutSecond * time.Second,
		MaxActive:   redisConfig.MaxActiveConnection,
		Wait:        redisConfig.WaitForNewConnection,

		Dial: func() (conn redis.Conn, e error) {
			c, err := redis.Dial("tcp", redisHost, redis.DialDatabase(redisConfig.Database), redis.DialPassword(redisConfig.Password))
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}

func cleanupHook(redisPool *redis.Pool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		_ = redisPool.Close()
		os.Exit(0)
	}()
}

// interface implementation

func (cache *RedisCacheManagerImp) GetString(key string) (string, error) {
	conn := cache.pool.Get()
	defer conn.Close()

	data, err := redis.String(conn.Do("GET", key))
	return data, err
}

func (cache *RedisCacheManagerImp) SetString(key string, value interface{}) error {
	conn := cache.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	return err
}
