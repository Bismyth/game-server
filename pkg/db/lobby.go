package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const lobbyHashName = "lobby"
const lobbyUsersHashName = "lobbyUsers"

type LobbyUser struct {
	Name string `json:"name"`
	Host bool   `json:"host"`
}

type LobbyUserList map[uuid.UUID]*LobbyUser

func CreateLobby(lobbyId uuid.UUID, initialUser uuid.UUID) error {
	err := SetLobbyProperty(lobbyId, "gameType", "")
	if err != nil {
		return err
	}

	err = SetLobbyProperty(lobbyId, "host", initialUser)
	if err != nil {
		return err
	}

	err = SetLobbyProperty(lobbyId, "inGame", false)
	if err != nil {
		return err
	}

	err = SaveLobbyUser(lobbyId, initialUser)
	if err != nil {
		return err
	}

	return nil
}

func LobbyInGame(lobbyId uuid.UUID) (bool, error) {
	return GetLobbyProperty[bool](lobbyId, "inGame")
}

func GetLobbyUsers(lobbyId uuid.UUID) (LobbyUserList, error) {
	conn := getConn()
	ctx := context.Background()

	list := make(LobbyUserList)

	m, err := conn.HGetAll(ctx, i(lobbyUsersHashName, lobbyId)).Result()
	if err != nil {
		return list, err
	}

	hostId, err := GetLobbyProperty[uuid.UUID](lobbyId, "host")
	if err != nil {
		return list, err
	}

	//validate keys are uuids
	for k := range m {
		id, err := uuid.Parse(k)
		if err != nil {
			return list, fmt.Errorf("invalid userid in list")
		}

		var user LobbyUser

		name, err := GetUserName(id)
		if err != nil {
			return list, fmt.Errorf("failed to get name for user")
		}
		user.Name = name
		if id == hostId {
			user.Host = true
		}

		list[id] = &user
	}

	return list, nil
}

func GetLobbyUserIds(lobbyId uuid.UUID) ([]uuid.UUID, error) {
	conn := getConn()
	ctx := context.Background()

	idStrings, err := conn.HKeys(ctx, i(lobbyUsersHashName, lobbyId)).Result()
	if err != nil {
		return nil, err
	}

	return ParseUUIDList(idStrings)
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

	host, err := IsUserLobbyHost(lobbyId, userId)
	if err != nil {
		return err
	}

	err = conn.HDel(ctx, i(lobbyUsersHashName, lobbyId), userId.String()).Err()
	if err != nil {
		return err
	}

	if host {
		users, err := GetLobbyUserIds(lobbyId)
		if err != nil {
			return err
		}
		if len(users) > 0 {
			err = SetLobbyProperty(lobbyId, "host", users[0])
			if err != nil {
				return err
			}
		} else {
			err = conn.Del(ctx, i(lobbyHashName, lobbyId)).Err()
			if err != nil {
				return err
			}
		}
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

func IsUserLobbyHost(lobbyId uuid.UUID, userId uuid.UUID) (bool, error) {
	id, err := GetHashTableProperty[uuid.UUID](i(lobbyHashName, lobbyId), "host")
	if err != nil {
		return false, err
	}

	return id == userId, nil
}

func LobbyExists(lobbyId uuid.UUID) (bool, error) {
	conn := getConn()
	ctx := context.Background()

	r, err := conn.Exists(ctx, i(lobbyHashName, lobbyId)).Result()
	if err != nil {
		return false, err
	}

	return r == 1, nil
}

func GetLobbyProperties(lobbyId uuid.UUID, fields []string) (map[string]interface{}, error) {
	return GetHashTableProperties(i(lobbyHashName, lobbyId), fields)
}

func GetLobbyProperty[T any](lobbyId uuid.UUID, field string) (T, error) {
	return GetHashTableProperty[T](i(lobbyHashName, lobbyId), field)
}

func SetLobbyProperty[T any](lobbyId uuid.UUID, field string, value T) error {
	return SetHashTableProperty(i(lobbyHashName, lobbyId), field, value)
}
