package api

import (
	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
)

const pt_OLobbyChange OPacketType = "server_lobby_change"

// Create Lobby Event
const pt_ICreateLobby IPacketType = "client_create_lobby"

func createLobby(i HandlerInput) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	newLobby := db.Lobby{Id: id, Users: map[uuid.UUID]bool{i.UserId: true}}

	newLobby.Save()

	return nil
}

// Join lobby event
const pt_IJoinLobby IPacketType = "client_join_lobby"

func joinLobby(i HandlerInput) error {
	lobbyId, err := hp[uuid.UUID](i.Packet)
	if err != nil {
		return err
	}

	lobby, err := db.GetLobby(*lobbyId)
	if err != nil {
		return err
	}

	lobby.Users[i.UserId] = true

	lobby.Save()

	sendLobbyChange(i.C, lobby)

	return nil
}

func sendLobbyChange(c ClientInterface, lobby *db.Lobby) error {
	newMessage := mp(pt_OLobbyChange, lobby.Users)

	keys := make([]uuid.UUID, 0, len(lobby.Users))
	for u := range lobby.Users {
		keys = append(keys, u)
	}

	SendMany(c, keys, &newMessage)

	return nil
}

// Leave lobby event
const pt_ILeaveLobby IPacketType = "client_leave_lobby"

func leaveLobby(i HandlerInput) error {
	lobbyId, err := hp[uuid.UUID](i.Packet)
	if err != nil {
		return err
	}

	lobby, err := db.GetLobby(*lobbyId)
	if err != nil {
		return err
	}

	delete(lobby.Users, i.UserId)

	lobby.Save()

	sendLobbyChange(i.C, lobby)

	return nil
}
