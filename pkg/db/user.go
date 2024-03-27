package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var userHashName = "user"

func (c *User) String() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		fmt.Println("failed to marshal user to redis")
		return "undefined"
	}

	return string(bytes)
}

func (c *User) Save() error {
	conn := getConn()

	data := c.String()

	err := conn.HSet(context.Background(), userHashName, map[string]interface{}{c.Id.String(): data}).Err()

	return err
}

func GetUser(id uuid.UUID) (*User, error) {
	conn := getConn()
	data, err := conn.HGet(context.Background(), userHashName, id.String()).Bytes()
	if err != nil {
		return nil, err
	}

	returnUser := User{Id: id}

	err = json.Unmarshal(data, &returnUser)
	if err != nil {
		return nil, err
	}

	return &returnUser, nil
}

func GetAllUserIds() ([]uuid.UUID, error) {
	conn := getConn()

	idStrings, err := conn.HKeys(context.Background(), userHashName).Result()
	if err != nil {
		return []uuid.UUID{}, err
	}

	outputIds := []uuid.UUID{}
	idErrors := []error{}

	for _, idString := range idStrings {
		id, err := uuid.Parse(idString)
		if err != nil {
			idErrors = append(idErrors, err)
		}
		outputIds = append(outputIds, id)
	}

	return outputIds, errors.Join(idErrors...)
}
