package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

var lobbyHashName = "lobby"

type LobbyUserList map[uuid.UUID]string

func marshalLobbyUsers(users LobbyUserList) (string, error) {
	usersString, err := json.Marshal(users)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user map: %v", err)
	}
	return string(usersString), nil
}

func CreateLobby(lobbyId uuid.UUID, initialUser uuid.UUID) error {
	conn := getConn()

	name, err := GetUserName(initialUser)
	if err != nil {
		return err
	}

	users := LobbyUserList{initialUser: name}
	usersString, err := marshalLobbyUsers(users)
	if err != nil {
		return err
	}

	return conn.HSet(context.Background(), i(lobbyHashName, lobbyId), "users", usersString).Err()
}

func GetLobbyUsers(lobbyId uuid.UUID) (LobbyUserList, error) {
	conn := getConn()
	ctx := context.Background()

	rawUsers, err := conn.HGet(ctx, i(lobbyHashName, lobbyId), "users").Bytes()
	var users LobbyUserList
	err = json.Unmarshal(rawUsers, &users)
	if err != nil {
		return users, err
	}

	return users, nil
}

func setLobbyUsers(lobbyId uuid.UUID, users LobbyUserList) error {
	conn := getConn()
	ctx := context.Background()

	usersString, err := marshalLobbyUsers(users)
	if err != nil {
		return err
	}

	return conn.HSet(ctx, i(lobbyHashName, lobbyId), "users", usersString).Err()
}

func AddLobbyUser(lobbyId uuid.UUID, userId uuid.UUID) error {
	users, err := GetLobbyUsers(lobbyId)
	if err != nil {
		return err
	}

	name, err := GetUserName(userId)
	if err != nil {
		return err
	}

	users[userId] = name
	return setLobbyUsers(lobbyId, users)
}

func RemoveLobbyUser(lobbyId uuid.UUID, userId uuid.UUID) error {
	users, err := GetLobbyUsers(lobbyId)
	if err != nil {
		return err
	}

	delete(users, userId)

	return setLobbyUsers(lobbyId, users)
}
