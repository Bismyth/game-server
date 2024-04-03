package db

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
)

const userHashName = "user"
const userLobbiesHashName = "userLobbies"

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

func SaveUserLobby(userId uuid.UUID, lobbyId uuid.UUID) error {
	conn := getConn()
	ctx := context.Background()

	err := conn.SAdd(ctx, i(userLobbiesHashName, userId), lobbyId.String()).Err()
	if err != nil {
		return err
	}

	return nil
}

func RemoveUserLobby(userId uuid.UUID, lobbyId uuid.UUID) error {
	conn := getConn()
	ctx := context.Background()

	err := conn.SRem(ctx, i(userLobbiesHashName, userId), lobbyId.String()).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetUserLobbies(userId uuid.UUID) ([]uuid.UUID, error) {
	conn := getConn()
	ctx := context.Background()

	idStrings, err := conn.SMembers(ctx, i(userLobbiesHashName, userId)).Result()
	if err != nil {
		return nil, err
	}
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
