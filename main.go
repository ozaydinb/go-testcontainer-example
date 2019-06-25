package main

import (
	"fmt"
	"github.com/ozaydinb/go-testcontainer-example/model"
	"github.com/ozaydinb/go-testcontainer-example/redis"
)

func main() {
	redisConfig := model.RedisConfig{
		Host:                 "127.0.0.1",
		Port:                 6379,
		WaitForNewConnection: true,
		MaxActiveConnection:  3,
		MaxIdle:              3,
		IdleTimeoutSecond:    5,
		Database:             0,
		Password:             "",
	}
	redisManager := redis.NewRedisCacheManager(redisConfig)

	err := redisManager.SetString("key1", "this is key1 value")

	if err != nil {
		fmt.Println(err.Error())
	}

	val, err := redisManager.GetString("key1")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(val)
}
