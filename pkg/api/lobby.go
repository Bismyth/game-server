package api

import (
	"log"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
)

const pt_OLobbyChange OPacketType = "server_lobby_change"
const pt_OJoinLobby OPacketType = "server_lobby_join"

// Create Lobby Event
const pt_ICreateLobby IPacketType = "client_lobby_create"

func createLobby(i HandlerInput) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	err = db.CreateLobby(id, i.UserId)
	if err != nil {
		return err
	}

	packet := mp(pt_OJoinLobby, id)
	Send(i.C, i.UserId, &packet)

	return nil
}

// Join lobby event
const pt_IJoinLobby IPacketType = "client_lobby_join"

func joinLobby(i HandlerInput) error {
	lobbyId, err := hp[uuid.UUID](i.Packet)
	if err != nil {
		return err
	}

	db.AddLobbyUser(*lobbyId, i.UserId)

	packet := mp(pt_OJoinLobby, *lobbyId)
	Send(i.C, i.UserId, &packet)

	sendLobbyChange(i.C, *lobbyId)

	return nil
}

// User request all lobby users
const pt_ILobbyUsers IPacketType = "client_lobby_users"

func lobbyUsers(i HandlerInput) error {
	lobbyId, err := hp[uuid.UUID](i.Packet)
	if err != nil {
		return err
	}

	_, packet := makeLobbyChangePacket(*lobbyId)
	Send(i.C, i.UserId, packet)

	return nil
}

func makeLobbyChangePacket(lobbyId uuid.UUID) ([]uuid.UUID, *Packet[[]string]) {
	users, err := db.GetLobbyUsers(lobbyId)
	if err != nil {
		log.Printf("Failed to retrieve list of users in lobby")
	}
	keys := make([]uuid.UUID, 0, len(users))
	names := make([]string, 0, len(users))
	for u, name := range users {
		keys = append(keys, u)
		names = append(names, name)
	}

	newMessage := mp(pt_OLobbyChange, names)
	return keys, &newMessage
}

func sendLobbyChange(c ClientInterface, lobbyId uuid.UUID) {
	ids, packet := makeLobbyChangePacket(lobbyId)
	SendMany(c, ids, packet)
}

// Leave lobby event
const pt_ILeaveLobby IPacketType = "client_lobby_leave"

func leaveLobby(i HandlerInput) error {
	lobbyId, err := hp[uuid.UUID](i.Packet)
	if err != nil {
		return err
	}

	err = db.RemoveLobbyUser(*lobbyId, i.UserId)
	if err != nil {
		return err
	}

	sendLobbyChange(i.C, *lobbyId)

	return nil
}
