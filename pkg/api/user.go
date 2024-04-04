package api

import (
	"log"
	"time"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
)

type m_User struct {
	Id       uuid.UUID   `json:"id"`
	Name     string      `json:"name"`
	LobbyIds []uuid.UUID `json:"lobbies"`
}

var nameGenerator = namegenerator.NewNameGenerator(time.Now().UTC().UnixNano())

func SetUserId(requestID uuid.UUID) uuid.UUID {
	id := requestID

	if !db.UserExists(id) {
		newId, err := uuid.NewV7()
		if err != nil {
			log.Println("Failed to generate new userId")
			return uuid.Nil
		}

		err = db.MakeUser(newId, nameGenerator.Generate())
		id = newId
		if err != nil {
			log.Println("Failed to make user in redis")
			id = uuid.Nil
		}
	}

	return id
}

// User init event
const pt_OUserInit OPacketType = "server_user_init"

func UserInitPacket(userId uuid.UUID) []byte {
	userMessage, err := makeUserMessage(userId)
	if err != nil {
		log.Println("Could not make user correctly")
	}

	oPacket := mp(pt_OUserInit, userMessage)

	return MarshalPacket(&oPacket)
}

func makeUserMessage(userId uuid.UUID) (m_User, error) {
	var user m_User
	user.Id = userId

	name, err := db.GetUserName(userId)
	if err != nil {
		return user, err
	}
	user.Name = name

	lobbies, err := db.GetUserLobbies(userId)
	if err != nil {
		return user, err
	}
	user.LobbyIds = lobbies

	return user, nil
}

const pt_OUserChange OPacketType = "server_user_change"

type m_UserChange struct {
	Name string
}

func sendUserChange(c Client, userId uuid.UUID) error {
	userMessage, err := makeUserMessage(userId)
	if err != nil {
		return err
	}

	oPacket := mp(pt_OUserChange, userMessage)
	Send(c, userId, &oPacket)
	return nil
}

// User change event
const pt_IUserNameChange IPacketType = "client_user_name_change"

func handleUserNameChange(i HandlerInput) error {
	iUserChange, err := hp[m_UserChange](i.Packet)
	if err != nil {
		return err
	}

	err = db.SetUserName(i.UserId, iUserChange.Name)
	if err != nil {
		return err
	}

	err = sendUserChange(i.C, i.UserId)
	if err != nil {
		return err
	}

	ids, err := db.GetUserLobbies(i.UserId)
	if err != nil {
		return err
	}

	for _, lobbyId := range ids {
		sendLobbyUserChange(i.C, lobbyId)
	}

	return nil
}
