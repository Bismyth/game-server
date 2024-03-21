package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

var config *Config = nil

func SetConfig(c *Config) {
	config = c
}

var conn *redis.Client

func getConn() *redis.Client {
	if conn == nil {
		conn = redis.NewClient(&redis.Options{
			Addr:     config.Addr,
			Password: config.Password,
			DB:       config.DB,
		})
	}

	response := conn.Ping(context.Background())
	if err := response.Err(); err != nil {
		panic(err)
	}

	return conn
}
