package db

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
)

const userHashName = "user"

func MakeUser(id uuid.UUID, name string) error {
	return SetUserName(id, name)
}

func UserExists(id uuid.UUID) bool {
	conn := getConn()

	i, _ := conn.Exists(context.Background(), i(userHashName, id)).Result()

	return i == 1
}

func GetUserName(id uuid.UUID) (string, error) {
	return GetHashTableProperty[string](i(userHashName, id), "name")
}

func SetUserName(id uuid.UUID, name string) error {
	return SetHashTableProperty(i(userHashName, id), "name", name)
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

	err := conn.SAdd(ctx, it(userHashName, userId, "lobbies"), lobbyId.String()).Err()
	if err != nil {
		return err
	}

	return nil
}

func RemoveUserLobby(userId uuid.UUID, lobbyId uuid.UUID) error {
	conn := getConn()
	ctx := context.Background()

	err := conn.SRem(ctx, it(userHashName, userId, "lobbies"), lobbyId.String()).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetUserLobbies(userId uuid.UUID) ([]uuid.UUID, error) {
	conn := getConn()
	ctx := context.Background()

	idStrings, err := conn.SMembers(ctx, it(userHashName, userId, "lobbies")).Result()
	if err != nil {
		return nil, err
	}

	return ParseUUIDList(idStrings)
}
