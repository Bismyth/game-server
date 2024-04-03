package db

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string `validate:"required" env:"ADDR"`
	Password string `env:"PASSWORD"`
	DB       int    `env:"DB"`
}

var config *Config = nil

func SetConfig(c Config) {
	config = &c
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

func i(namespace string, id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", namespace, id.String())
}

func ia(namespace string) string {
	return fmt.Sprintf("%s:*", namespace)
}
