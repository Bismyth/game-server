package ws

import (
	"encoding/json"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
)

func (h *Hub) createLobby(message *BaseMessage) error {

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	newLobby := db.Lobby{Id: id, Users: map[uuid.UUID]bool{message.client.id: true}}

	newLobby.Save()

	return nil
}

func (h *Hub) joinLobby(message *BaseMessage) error {

	var lobbyId uuid.UUID
	err := json.Unmarshal(message.Data, &lobbyId)
	if err != nil {
		return err
	}

	lobby, err := db.GetLobby(lobbyId)
	if err != nil {
		return err
	}

	lobby.Users[message.client.id] = true

	lobby.Save()

	h.broadcastCurrentUsers(lobby)

	return nil
}

func (h *Hub) broadcastCurrentUsers(lobby *db.Lobby) error {

	rawData, err := json.Marshal(lobby.Users)
	if err != nil {
		return err
	}

	newMessage := BaseMessage{
		Type: "lobbyUsers",
		Data: rawData,
	}

	for id := range lobby.Users {
		h.clientIds[id].send <- &newMessage
	}

	return nil
}

func (h *Hub) leaveLobby(message *BaseMessage) error {
	var lobbyId uuid.UUID
	err := json.Unmarshal(message.Data, &lobbyId)
	if err != nil {
		return err
	}

	lobby, err := db.GetLobby(lobbyId)
	if err != nil {
		return err
	}

	delete(lobby.Users, message.client.id)

	lobby.Save()

	h.broadcastCurrentUsers(lobby)

	return nil
}
