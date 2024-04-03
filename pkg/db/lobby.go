package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const lobbyHashName = "lobby"
const lobbyUsersHashName = "lobbyUsers"

type LobbyUserList map[uuid.UUID]string

func CreateLobby(lobbyId uuid.UUID, initialUser uuid.UUID) error {
	conn := getConn()
	ctx := context.Background()

	err := conn.HSet(ctx, i(lobbyHashName, lobbyId), "gameType", "").Err()
	if err != nil {
		return err
	}

	err = SaveLobbyUser(lobbyId, initialUser)
	if err != nil {
		return err
	}

	return nil
}

func GetLobbyUsers(lobbyId uuid.UUID) (LobbyUserList, error) {
	conn := getConn()
	ctx := context.Background()

	list := make(LobbyUserList)

	m, err := conn.HGetAll(ctx, i(lobbyUsersHashName, lobbyId)).Result()
	if err != nil {
		return list, err
	}

	//validate keys are uuids
	for k := range m {
		id, err := uuid.Parse(k)
		if err != nil {
			return list, fmt.Errorf("invalid userid in list")
		}

		name, err := GetUserName(id)
		if err != nil {
			return list, fmt.Errorf("failed to get name for user")
		}

		list[id] = name
	}

	return list, nil
}

func SaveLobbyUser(lobbyId uuid.UUID, userId uuid.UUID) error {
	conn := getConn()
	ctx := context.Background()

	err := conn.HSet(ctx, i(lobbyUsersHashName, lobbyId), userId.String(), "").Err()
	if err != nil {
		return err
	}

	err = SaveUserLobby(userId, lobbyId)
	if err != nil {
		return err
	}

	return nil
}

func RemoveLobbyUser(lobbyId uuid.UUID, userId uuid.UUID) error {
	conn := getConn()
	ctx := context.Background()

	err := conn.HDel(ctx, i(lobbyUsersHashName, lobbyId), userId.String()).Err()
	if err != nil {
		return err
	}

	err = RemoveUserLobby(userId, lobbyId)
	if err != nil {
		return err
	}

	return nil
}

func IsUserInLobby(lobbyId uuid.UUID, userId uuid.UUID) (bool, error) {
	conn := getConn()
	ctx := context.Background()

	r, err := conn.HExists(ctx, i(lobbyUsersHashName, lobbyId), userId.String()).Result()
	if err != nil {
		return false, err
	}

	return r, nil
}
