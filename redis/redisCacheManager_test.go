package redis

import (
	"context"
	"github.com/ozaydinb/go-testcontainer-example/model"
	"github.com/testcontainers/testcontainers-go"
	"gotest.tools/assert"
	"testing"
)

func TestRedisCacheManager_Should_Get_Value_When_Key_Exists(t *testing.T) {
	ip, redisPort := startRedisContainer(t)
	redisManager := getRedisCacheManager(ip, redisPort)

	err := redisManager.SetString("key1", "val1")
	if err != nil {
		t.Error(err)
	}

	val, err := redisManager.GetString("key1")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, val, "val1")
}

func TestRedisCacheManager_Should_Return_Error_When_Key_Not_Exists(t *testing.T) {
	ip, redisPort := startRedisContainer(t)
	redisManager := getRedisCacheManager(ip, redisPort)

	_, err := redisManager.GetString("key1")
	assert.ErrorContains(t, err, "redigo: nil returned")

}

func TestRedisCacheManager_Should_Set_String_Value(t *testing.T) {
	ip, redisPort := startRedisContainer(t)
	redisManager := getRedisCacheManager(ip, redisPort)

	err := redisManager.SetString("key1", "val1")
	assert.NilError(t, err)
}

func getRedisCacheManager(ip string, redisPort int) RedisCacheManager {
	redisConfig := model.RedisConfig{
		Host:                 ip,
		Password:             "",
		Database:             0,
		IdleTimeoutSecond:    3,
		MaxIdle:              3,
		MaxActiveConnection:  3,
		WaitForNewConnection: true,
		Port:                 redisPort,
	}
	redisManager := NewRedisCacheManager(redisConfig)
	return redisManager
}

func startRedisContainer(t *testing.T) (string, int) {
	ctx := context.Background()
	request := testcontainers.ContainerRequest{
		Image:        "redis",
		ExposedPorts: []string{"6379/tcp"},
	}

	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})

	if err != nil {
		t.Error(err)
	}

	defer func() {
		err := redisContainer.Terminate(ctx)
		if err != nil {
			t.Error(err)
		}
	}()

	ip, err := redisContainer.Host(ctx)

	if err != nil {
		t.Error(err)
	}
	redisPort, err := redisContainer.MappedPort(ctx, "6379/tcp")
	if err != nil {
		t.Error(err)
	}

	err = redisContainer.Start(ctx)
	if err != nil {
		t.Error(err)
	}

	return ip, redisPort.Int()

}
