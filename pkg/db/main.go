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

func it(namespace string, id uuid.UUID, field string) string {
	return fmt.Sprintf("%s:%s:%s", namespace, id.String(), field)
}

func GetHashTableProperties(key string, fields []string) (map[string]interface{}, error) {
	conn := getConn()
	ctx := context.Background()

	r, err := conn.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return nil, err
	}

	output := make(map[string]interface{})

	for i, field := range fields {
		output[field] = r[i]
	}

	return output, nil
}

func GetHashTableProperty(key, field string) (string, error) {
	conn := getConn()
	ctx := context.Background()

	r, err := conn.HGet(ctx, key, field).Result()
	if err != nil {
		return "", err
	}

	return r, nil
}
