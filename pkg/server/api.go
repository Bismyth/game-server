package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Bismyth/game-server/pkg/api"
	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
)

func (S *Server) RegisterAPI(mux *http.ServeMux) {
	mux.HandleFunc("/api/room/create", S.CreateRoom)
	mux.HandleFunc("/api/room/join", S.JoinRoom)
	mux.HandleFunc("/api/room/tokens", S.ValidateRoomTokens)
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

func (S *Server) ValidateRoomTokens(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post is allowed", http.StatusMethodNotAllowed)
		return
	}

	var input []string
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "expected a list of ids", http.StatusUnprocessableEntity)
		return
	}

	output := make([]bool, len(input))
	for i, idString := range input {
		id, err := uuid.Parse(idString)
		if err != nil {
			output[i] = false
			continue
		}

		exists, _ := db.RoomExists(id)
		output[i] = exists
	}

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		log.Printf("failed to marshal api response: %v", err)
	}
}
