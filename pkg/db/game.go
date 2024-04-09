package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const gameHashName = "game"

func SetGameCache(gameId uuid.UUID, data any) error {
	conn := getConn()
	ctx := context.Background()

	v, err := Encode(data)
	if err != nil {
		return err
	}

	err = conn.Set(ctx, it(gameHashName, gameId, "cache"), v, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetGameCache[T any](gameId uuid.UUID) (T, error) {
	conn := getConn()
	ctx := context.Background()

	output := new(T)

	r, err := conn.Get(ctx, it(gameHashName, gameId, "cache")).Bytes()
	if err != nil {
		return *output, err
	}

	err = Decode(r, output)
	if err != nil {
		return *output, err
	}

	return *output, nil
}

func ExpireCache(gameId uuid.UUID, duration time.Duration) error {
	conn := getConn()
	ctx := context.Background()

	err := conn.Expire(ctx, it(gameHashName, gameId, "cache"), duration).Err()
	if err != nil {
		return err
	}

	return nil
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

func GetMultiPlayerProperty[T any](gameId uuid.UUID, field string, playerType string) ([]T, error) {
	players, err := PlayerTypeGetAll(gameId, playerType)
	if err != nil {
		return nil, err
	}

	output := []T{}

	for _, player := range players {
		p, err := GetPlayerProperty[T](gameId, player, field)
		if err != nil {
			return nil, err
		}
		output = append(output, p)
	}

	return output, nil
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
