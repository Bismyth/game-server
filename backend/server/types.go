package server

import "github.com/google/uuid"

type RoomStatus string

const roomstatus_success = "success"
const roomstatus_error = "error"

type RoomResponse struct {
	Status RoomStatus `json:"status"`
	Error  string     `json:"error"`
	Token  string     `json:"token"`
	Id     uuid.UUID  `json:"id"`
}

type RoomJoinInput struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
