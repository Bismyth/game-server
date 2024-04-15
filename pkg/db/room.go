package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const roomHash = "room"

const roomNameHash = "name"
const roomSessionHash = "session"

func CreateRoom() (uuid.UUID, error) {
	roomId, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}

	err = SetRoomProperty(roomId, "gameType", "")
	if err != nil {
		return roomId, err
	}

	err = SetRoomProperty(roomId, "inGame", false)
	if err != nil {
		return roomId, err
	}

	return roomId, nil
}

func SetRoomHost(roomId uuid.UUID, userId uuid.UUID) error {
	err := SetRoomProperty(roomId, "host", userId)
	if err != nil {
		return err
	}

	return nil
}

func NewRoomUser(roomId uuid.UUID, name string) (uuid.UUID, error) {
	userId, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}

	err = SetRoomUserName(roomId, userId, name)
	if err != nil {
		return userId, err
	}

	return userId, nil
}

func RoomExists(roomId uuid.UUID) (bool, error) {
	conn := getConn()
	ctx := context.Background()

	i, err := conn.Exists(ctx, i(roomHash, roomId)).Result()
	if err != nil {
		return false, err
	}

	return i == 1, nil
}

func RemoveRoomUser(roomId uuid.UUID, userId uuid.UUID) error {
	conn := getConn()
	ctx := context.Background()

	err := conn.HDel(ctx, it(roomHash, roomId, roomNameHash), userId.String()).Err()
	if err != nil {
		return err
	}

	err = conn.HDel(ctx, it(roomHash, roomId, roomSessionHash), userId.String()).Err()
	if err != nil {
		return err
	}

	hostId, err := GetRoomProperty[uuid.UUID](roomId, "host")
	if err != nil {
		return err
	}

	if userId == hostId {
		users, err := GetRoomUsers(roomId)
		if err != nil {
			return err
		}
		if len(users) > 0 {
			var newHost uuid.UUID
			for id := range users {
				newHost = id
				break
			}

			err = SetRoomProperty(roomId, "host", newHost)
			if err != nil {
				return err
			}
		} else {
			err = conn.Del(ctx, i(roomHash, roomId)).Err()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func RemoveRoomUserSession(roomId uuid.UUID, userId uuid.UUID) error {
	conn := getConn()
	ctx := context.Background()

	err := conn.HDel(ctx, it(roomHash, roomId, roomSessionHash), userId.String()).Err()
	if err != nil {
		return err
	}

	return nil
}

func RoomHasUser(roomId uuid.UUID, userId uuid.UUID) (bool, error) {
	conn := getConn()
	ctx := context.Background()

	exists, err := conn.HExists(ctx, it(roomHash, roomId, roomNameHash), userId.String()).Result()
	if err != nil {
		return false, err
	}

	return exists, nil
}

func SetRoomUserName(roomId uuid.UUID, userId uuid.UUID, name string) error {
	return SetHashTableProperty(it(roomHash, roomId, roomNameHash), userId.String(), name)
}

func SetRoomUserSession(roomId uuid.UUID, userId uuid.UUID, sessionId uuid.UUID) error {
	return SetHashTableProperty(it(roomHash, roomId, roomSessionHash), userId.String(), sessionId)
}

func GetRoomUserSession(roomId uuid.UUID, userId uuid.UUID) (uuid.UUID, error) {
	return GetHashTableProperty[uuid.UUID](it(roomHash, roomId, roomSessionHash), userId.String())
}

func SetRoomProperty[T any](roomId uuid.UUID, field string, data T) error {
	return SetHashTableProperty(i(roomHash, roomId), field, data)
}

func GetRoomProperty[T any](roomId uuid.UUID, field string) (T, error) {
	return GetHashTableProperty[T](i(roomHash, roomId), field)
}

func GetRoomProperties(roomId uuid.UUID, fields []string) (map[string]interface{}, error) {
	return GetHashTableProperties(i(roomHash, roomId), fields)
}

func GetRoomSessions(roomId uuid.UUID) ([]uuid.UUID, error) {
	conn := getConn()
	ctx := context.Background()

	idStrings, err := conn.HVals(ctx, it(roomHash, roomId, roomSessionHash)).Result()
	if err != nil {
		return nil, err
	}

	return ParseUUIDList(idStrings)
}

func GetRoomUserName(roomId uuid.UUID, userId uuid.UUID) (string, error) {
	return GetHashTableProperty[string](it(roomHash, roomId, roomNameHash), userId.String())
}

type UserList map[uuid.UUID]string

func GetRoomUsers(roomId uuid.UUID) (UserList, error) {
	conn := getConn()
	ctx := context.Background()

	rawMap, err := conn.HGetAll(ctx, it(roomHash, roomId, roomNameHash)).Result()
	if err != nil {
		return nil, err
	}

	output := make(UserList)
	for idString, name := range rawMap {
		id, err := uuid.Parse(idString)
		if err != nil {
			return nil, fmt.Errorf("invalid uuid in user list")
		}
		output[id] = name
	}

	return output, nil
}

func GetRoomUserOrder(roomId uuid.UUID) ([]uuid.UUID, error) {
	conn := getConn()
	ctx := context.Background()

	idStrings, err := conn.HKeys(ctx, it(roomHash, roomId, roomNameHash)).Result()
	if err != nil {
		return nil, err
	}

	return ParseUUIDList(idStrings)
}
