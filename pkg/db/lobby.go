package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

var lobbyHashName = "lobby"

func (c *Lobby) String() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		fmt.Println("failed to marshal lobby to redis")
		return "undefined"
	}

	return string(bytes)
}

func (l *Lobby) Save() error {
	conn := getConn()

	data := l.String()

	err := conn.HSet(context.Background(), lobbyHashName, map[string]interface{}{l.Id.String(): data}).Err()

	return err
}

func GetLobby(id uuid.UUID) (*Lobby, error) {
	conn := getConn()
	data, err := conn.HGet(context.Background(), lobbyHashName, id.String()).Bytes()
	if err != nil {
		return nil, err
	}

	returnLobby := Lobby{Id: id}

	err = json.Unmarshal(data, &returnLobby)
	if err != nil {
		return nil, err
	}

	return &returnLobby, nil
}
