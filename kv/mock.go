package kv

import (
	"context"
	"strconv"

	"github.com/alicebob/miniredis/v2"
)

var redisServer *miniredis.Miniredis

func InitMock(ctx context.Context, config *Config) (*Client, error) {

	if config == nil {
		config = &Config{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		}
	}

	var err error
	redisServer, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	serverPort, _ := strconv.Atoi(redisServer.Port())
	config.Host = redisServer.Host()
	config.Port = serverPort
	return Init(ctx, config)
}
