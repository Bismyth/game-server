package db

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
)

var userHashName = "user"

func MakeUser(id uuid.UUID, name string) error {
	conn := getConn()

	return conn.HSet(context.Background(), i(userHashName, id), "name", name).Err()
}

func UserExists(id uuid.UUID) bool {
	conn := getConn()

	i, _ := conn.Exists(context.Background(), i(userHashName, id)).Result()

	return i == 1
}

func GetUserName(id uuid.UUID) (string, error) {
	conn := getConn()
	name, err := conn.HGet(context.Background(), i(userHashName, id), "name").Result()
	if err != nil {
		return "", err
	}

	return name, nil
}

func SetUserName(id uuid.UUID, name string) error {
	conn := getConn()

	return conn.HSet(context.Background(), i(userHashName, id), "name", name).Err()
}

func GetAllUserIds() ([]uuid.UUID, error) {
	conn := getConn()

	ids := []uuid.UUID{}
	idErrors := []error{}
	ctx := context.Background()
	iter := conn.Scan(ctx, 0, ia(userHashName), 0).Iterator()

	for iter.Next(ctx) {
		id, err := uuid.Parse(iter.Val())
		if err != nil {
			idErrors = append(idErrors, err)
			log.Println("Invalid uuid found in all user search")
			continue
		}
		ids = append(ids, id)
	}

	return ids, errors.Join(idErrors...)
}

func SetUserLobby(id uuid.UUID, lobbyId uuid.UUID) error {
	conn := getConn()
	ctx := context.Background()

	err := conn.HDel(ctx, i(userHashName, id), "game").Err()
	if err != nil {
		return err
	}

	err = conn.HSet(ctx, i(userHashName, id), "lobby", lobbyId.String()).Err()
	if err != nil {
		return err
	}

	return nil
}

func SetUserGame(id uuid.UUID, gameId uuid.UUID) error {
	conn := getConn()
	ctx := context.Background()

	err := conn.HDel(ctx, i(userHashName, id), "lobby").Err()
	if err != nil {
		return err
	}

	err = conn.HSet(ctx, i(userHashName, id), "game", gameId.String()).Err()
	if err != nil {
		return err
	}

	return nil
}

func IsUserInGame(userId uuid.UUID, gameId uuid.UUID) bool {
	conn := getConn()
	ctx := context.Background()

	idString, err := conn.HGet(ctx, i(userHashName, userId), "game").Result()
	if err != nil {
		return false
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return false
	}

	return id == gameId
}
