package db

import (
	"context"

	"github.com/google/uuid"
)

const sessionHash = "session"

func NewSession(sessionId, roomId, userId uuid.UUID) error {
	err := SetHashTableProperty(i(sessionHash, sessionId), "room", roomId)
	if err != nil {
		return err
	}

	err = SetHashTableProperty(i(sessionHash, sessionId), "user", userId)
	if err != nil {
		return err
	}

	return nil
}

type Session struct {
	RoomId uuid.UUID
	UserId uuid.UUID
}

func GetSessionDetails(sessionId uuid.UUID) (*Session, error) {
	roomId, err := GetHashTableProperty[uuid.UUID](i(sessionHash, sessionId), "room")
	if err != nil {
		return nil, err
	}

	userId, err := GetHashTableProperty[uuid.UUID](i(sessionHash, sessionId), "user")
	if err != nil {
		return nil, err
	}

	return &Session{
		RoomId: roomId,
		UserId: userId,
	}, nil
}

func DeleteSession(sessionId uuid.UUID) error {
	conn := getConn()
	ctx := context.Background()

	err := conn.Del(ctx, i(sessionHash, sessionId)).Err()
	if err != nil {
		return err
	}

	return nil
}
