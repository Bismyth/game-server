package db

import "github.com/google/uuid"

type User struct {
	Id   uuid.UUID `json:"-"`
	Name string    `json:"name"`
}

type Lobby struct {
	Id uuid.UUID `json:"-"`
	// Private bool
	// Hashed link?
	// Name
	Users map[uuid.UUID]bool
}
