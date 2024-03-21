package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

var clientHashName = "client"

func (c *User) String() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		fmt.Println("failed to marshal client to redis")
		return "undefined"
	}

	return string(bytes)
}

func (c *User) Save() error {
	conn := getConn()

	data := c.String()

	err := conn.HSet(context.Background(), clientHashName, map[string]interface{}{c.Id.String(): data}).Err()

	return err
}

func GetClient(id uuid.UUID) (*User, error) {
	conn := getConn()
	data, err := conn.HGet(context.Background(), clientHashName, id.String()).Bytes()
	if err != nil {
		return nil, err
	}

	returnClient := User{Id: id}

	err = json.Unmarshal(data, &returnClient)
	if err != nil {
		return nil, err
	}

	return &returnClient, nil
}
