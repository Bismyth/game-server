package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"

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

func ic(key string) string {
	return fmt.Sprintf("%s:cursor", key)
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

func GetHashTableProperty[T any](key, field string) (T, error) {
	conn := getConn()
	ctx := context.Background()

	result := new(T)

	r, err := conn.HGet(ctx, key, field).Bytes()
	if err != nil {
		return *result, err
	}

	err = Decode(r, result)
	if err != nil {
		return *result, err
	}

	return *result, nil
}

func SetHashTableProperty[T any](key, field string, value T) error {
	conn := getConn()
	ctx := context.Background()

	stringValue, err := Encode(value)
	if err != nil {
		return err
	}

	err = conn.HSet(ctx, key, field, stringValue).Err()
	if err != nil {
		return err
	}

	return nil
}

func ParseUUIDList(idStrings []string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, len(idStrings))
	for i, idString := range idStrings {
		id, err := uuid.Parse(idString)
		if err != nil {
			return ids, err
		}

		ids[i] = id
	}
	return ids, nil
}

func Encode(t interface{}) (string, error) {
	switch v := t.(type) {
	case string:
		return v, nil
	case *string:
		return *v, nil
	case int:
		return strconv.Itoa(v), nil
	case *int:
		return strconv.Itoa(*v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case *int64:
		return strconv.FormatInt(*v, 10), nil
	case uuid.UUID:
		return v.String(), nil
	case *uuid.UUID:
		return v.String(), nil
	case bool:
		return strconv.FormatBool(v), nil
	case *bool:
		return strconv.FormatBool(*v), nil
	default:
		raw, err := json.Marshal(t)
		return string(raw), err
	}
}

func Decode(raw []byte, output interface{}) error {
	rv := reflect.ValueOf(output)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("not a pointer")
	}

	switch t := output.(type) {
	case *string:
		s := string(raw)
		*t = s
	case *int:
		s := string(raw)
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*t = i
	case *int64:
		i, err := strconv.ParseInt(string(raw), 10, 0)
		if err != nil {
			return err
		}
		*t = i
	case *uuid.UUID:
		id, err := uuid.Parse(string(raw))
		if err != nil {
			return err
		}
		*t = id
	case *bool:
		b, err := strconv.ParseBool(string(raw))
		if err != nil {
			return err
		}
		*t = b
	default:
		return json.Unmarshal(raw, output)
	}

	return nil
}

func IsNilErr(err error) bool {
	return errors.Is(err, redis.Nil)
}
