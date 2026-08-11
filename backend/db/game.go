package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const gameHashName = "game"

func SetGameKey(gameId uuid.UUID, key string, data any) error {
	conn := getConn()
	ctx := context.Background()

	v, err := Encode(data)
	if err != nil {
		return err
	}

	err = conn.Set(ctx, it(gameHashName, gameId, key), v, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func SetGameCache(gameId uuid.UUID, data any) error {
	return SetGameKey(gameId, "cache", data)
}

func GetGameKey[T any](gameId uuid.UUID, key string) (T, error) {
	conn := getConn()
	ctx := context.Background()

	output := new(T)
	v, err := conn.Exists(ctx, it(gameHashName, gameId, key)).Result()
	if err != nil {
		return *output, err
	}
	if v <= 0 {
		return *output, nil
	}

	r, err := conn.Get(ctx, it(gameHashName, gameId, key)).Bytes()
	if err != nil {
		return *output, err
	}

	err = Decode(r, output)
	if err != nil {
		return *output, err
	}

	return *output, nil
}

func GetGameCache[T any](gameId uuid.UUID) (T, error) {
	return GetGameKey[T](gameId, "cache")
}

func ExpireGameKey(gameId uuid.UUID, key string, duration time.Duration) error {
	conn := getConn()
	ctx := context.Background()

	err := conn.Expire(ctx, it(gameHashName, gameId, key), duration).Err()
	if err != nil {
		return err
	}

	return nil
}

func ExpireCache(gameId uuid.UUID, duration time.Duration) error {
	return ExpireGameKey(gameId, "cache", duration)
}

func SetGameProperty(gameId uuid.UUID, field string, data any) error {
	return SetHashTableProperty(i(gameHashName, gameId), field, data)
}

func SetPlayerProperty(gameId uuid.UUID, playerId uuid.UUID, field string, data any) error {
	return SetGameProperty(gameId, i(field, playerId), data)
}

func GetGameProperty[T any](gameId uuid.UUID, field string) (T, error) {
	return GetHashTableProperty[T](i(gameHashName, gameId), field)
}

func GetPlayerProperty[T any](gameId uuid.UUID, playerId uuid.UUID, field string) (T, error) {
	return GetGameProperty[T](gameId, i(field, playerId))
}

func DeleteGame(gameId uuid.UUID) error {
	conn := getConn()
	ctx := context.Background()

	err := conn.Del(ctx, i(gameHashName, gameId)).Err()
	if err != nil {
		return err
	}

	return nil
}
