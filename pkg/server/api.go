package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Bismyth/game-server/pkg/api"
	"github.com/google/uuid"
)

func (S *Server) RegisterAPI(mux *http.ServeMux) {
	mux.HandleFunc("/api/room/create", S.CreateRoom)
	mux.HandleFunc("/api/room/join", S.JoinRoom)
}

func (S *Server) CreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post is allowed", http.StatusMethodNotAllowed)
		return
	}

	var input RoomJoinInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Name == "" {
		http.Error(w, "bad input", http.StatusUnprocessableEntity)
		return
	}

	roomId, token, err := api.CreateRoom(input.Name)
	sendRoomResponse(w, roomId, token, err)
}

func sendRoomResponse(w http.ResponseWriter, roomId uuid.UUID, token string, err error) {
	var response RoomResponse

	responseWriter := json.NewEncoder(w)
	if err != nil {
		response = RoomResponse{
			Status: roomstatus_error,
			Error:  err.Error(),
			Id:     roomId,
		}
	} else {
		response = RoomResponse{
			Status: roomstatus_success,
			Token:  token,
			Id:     roomId,
		}
	}
	err = responseWriter.Encode(response)
	if err != nil {
		log.Printf("failed to marshal api response: %v", err)
	}
}

func (S *Server) JoinRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post is allowed", http.StatusMethodNotAllowed)
		return
	}

	var input RoomJoinInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Name == "" {
		http.Error(w, "bad input", http.StatusUnprocessableEntity)
		return
	}

	token, err := api.JoinRoom(S.WSHub, input.Id, input.Name)
	sendRoomResponse(w, input.Id, token, err)
}
