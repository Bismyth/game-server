package api

import (
	"fmt"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
)

// HTTP POST Path for creating Room and getting token
func CreateRoom(name string) (uuid.UUID, string, error) {
	roomId, err := db.CreateRoom()
	if err != nil {
		return uuid.Nil, "", err
	}

	userId, err := db.NewRoomUser(roomId, name)
	if err != nil {
		return roomId, "", err
	}

	err = db.SetRoomHost(roomId, userId)
	if err != nil {
		return roomId, "", err
	}

	token, err := GenerateRoomToken(roomId, userId)
	if err != nil {
		return roomId, "", err
	}

	return roomId, token, nil
}

// HTTP Post Path for joining room and getting room token
func JoinRoom(c Client, roomId uuid.UUID, name string) (string, error) {
	exists, err := db.RoomExists(roomId)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("room does not exist")
	}

	userId, err := db.NewRoomUser(roomId, name)
	if err != nil {
		return "", err
	}
	token, err := GenerateRoomToken(roomId, userId)
	if err != nil {
		return "", err
	}

	sendRoomUsers(c, roomId)

	return token, nil
}
