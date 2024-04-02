package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

const gameHashName = "game"

func NewGame(gameType string, players []uuid.UUID) (uuid.UUID, error) {
	conn := getConn()
	ctx := context.Background()

	newGameId, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create game id")
	}

	rawPlayers, err := json.Marshal(players)
	if err != nil {
		return newGameId, fmt.Errorf("failed to marshal player ids")
	}

	err = conn.HSet(ctx, i(gameHashName, newGameId), map[string]string{
		"type":     gameType,
		"players":  string(rawPlayers),
		"#players": strconv.Itoa(len(players)),
	}).Err()
	if err != nil {
		return newGameId, fmt.Errorf("failed to create game")
	}

	return newGameId, nil
}

func GetGameType(id uuid.UUID) (string, error) {
	conn := getConn()
	ctx := context.Background()

	gameType, err := conn.HGet(ctx, i(gameHashName, id), "type").Result()
	if err != nil {
		return "", fmt.Errorf("could not get game type")
	}

	return gameType, nil
}

func GetGamePlayers(id uuid.UUID) ([]uuid.UUID, error) {
	conn := getConn()
	ctx := context.Background()

	var players []uuid.UUID

	rawPlayers, err := conn.HGet(ctx, i(gameHashName, id), "players").Bytes()
	if err != nil {
		return players, fmt.Errorf("could not get game type")
	}

	err = json.Unmarshal(rawPlayers, &players)
	if err != nil {
		return players, fmt.Errorf("could not unmarshal player ids")
	}

	return players, nil
}

func SetGameProperty(gameId uuid.UUID, field string, data any) error {
	conn := getConn()
	ctx := context.Background()

	rawData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal %s data", field)
	}

	err = conn.HSet(ctx, i(gameHashName, gameId), field, string(rawData)).Err()
	if err != nil {
		return fmt.Errorf("failed to set property: %s", field)
	}

	return nil
}

func SetPlayerProperty(gameId uuid.UUID, playerId uuid.UUID, field string, data any) error {
	return SetGameProperty(gameId, i(field, playerId), data)
}

func GetGameProperty[T any](gameId uuid.UUID, field string) (T, error) {
	conn := getConn()
	ctx := context.Background()

	var data T

	rawData, err := conn.HGet(ctx, i(gameHashName, gameId), field).Bytes()
	if err != nil {
		return data, fmt.Errorf("failed to get property: %s", field)
	}

	err = json.Unmarshal(rawData, &data)
	if err != nil {
		return data, fmt.Errorf("failed to unmarshal data for: %s", field)
	}

	return data, nil
}

func GetPlayerProperty[T any](gameId uuid.UUID, playerId uuid.UUID, field string) (T, error) {
	return GetGameProperty[T](gameId, i(field, playerId))
}

func GetMultiPlayerProperty[T any](gameId uuid.UUID, field string) ([]T, error) {
	players, err := GetGamePlayers(gameId)
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
