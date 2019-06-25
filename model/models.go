package model

import "time"

type RedisConfig struct {
	MaxIdle int
	IdleTimeoutSecond time.Duration
	MaxActiveConnection int
	WaitForNewConnection bool
	Database int
	Password string
	Host string
	Port int

}
