package db

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string `validate:"required"`
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
		log.Panicf("failed to get db connection: %v", err)
	}

	return conn
}

// Deletes all data in the redis database
func FlushDB() error {
	db := getConn()
	return db.FlushDB(context.Background()).Err()
}
