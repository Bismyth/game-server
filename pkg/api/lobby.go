package api

import (
	"log"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

// All outgoing lobby data
type m_Lobby struct {
	Id          uuid.UUID
	MaxPlayers  int
	MinPlayers  int
	Name        string
	GameType    string
	GameOptions interface{}
}

// incoming client lobby change
type m_LobbyChange struct {
	MaxPlayers int
	MinPlayers int
	Name       string
	GameType   string
}

const pt_OLobbyChange OPacketType = "server_lobby_change"

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

	return nil
}

// Join lobby event
const pt_IJoinLobby IPacketType = "client_lobby_join"

func joinLobby(i HandlerInput) error {
	lobbyId, err := hp[uuid.UUID](i.Packet)
	if err != nil {
		return err
	}

	db.SaveLobbyUser(*lobbyId, i.UserId)

	sendLobbyUserChange(i.C, *lobbyId)

	return nil
}

// User request all lobby users
const pt_ILobbyUsers IPacketType = "client_lobby_users"

func lobbyUsers(i HandlerInput) error {
	lobbyId, err := hp[uuid.UUID](i.Packet)
	if err != nil {
		return err
	}

	_, packet := makeLobbyUsersMessage(*lobbyId)
	Send(i.C, i.UserId, packet)

	return nil
}

func makeLobbyUsersMessage(lobbyId uuid.UUID) ([]uuid.UUID, *Packet[db.LobbyUserList]) {
	users, err := db.GetLobbyUsers(lobbyId)
	if err != nil {
		log.Printf("Failed to retrieve list of users in lobby")
	}
	keys := make([]uuid.UUID, 0, len(users))

	for u := range users {
		keys = append(keys, u)
	}

	newMessage := mp(pt_OLobbyChange, users)
	return keys, &newMessage
}

func sendLobbyUserChange(c interfaces.Client, lobbyId uuid.UUID) {
	ids, packet := makeLobbyUsersMessage(lobbyId)
	SendMany(c, ids, packet)
}

func sendToLobbyUser(c interfaces.Client, lobbyId uuid.UUID, userId uuid.UUID) {

}

func sendToLobbyUsers(c interfaces.Client, lobbyId uuid.UUID) {

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

	sendLobbyUserChange(i.C, *lobbyId)

	return nil
}
