package api

import (
	"log"
	"time"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
)

type m_User struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

var nameGenerator = namegenerator.NewNameGenerator(time.Now().UTC().UnixNano())

func SetUserId(requestID uuid.UUID) uuid.UUID {
	user, err := db.GetUser(requestID)

	if err != nil {
		newId, err := uuid.NewV7()
		if err != nil {
			log.Println("Failed to generate new userId")
			return uuid.Nil
		}

		user = &db.User{Id: newId, Name: nameGenerator.Generate()}
	}

	user.Save()

	return user.Id
}

// User init event
const pt_OUserInit OPacketType = "server_user_init"

func SendUserInit(c ClientInterface, userId uuid.UUID) error {
	user, err := db.GetUser(userId)
	if err != nil {
		log.Println("Failed to initilize user")
	}

	oData := m_User{
		Id:   user.Id,
		Name: user.Name,
	}

	oPacket := mp(pt_OUserInit, oData)
	Send(c, user.Id, &oPacket)

	return nil
}

// User change event
const pt_IUserChange IPacketType = "client_user_change"
const pt_OUserChange OPacketType = "server_user_change"

func handleUserChange(i HandlerInput) error {
	iUserChange, err := hp[m_User](i.Packet)
	if err != nil {
		return err
	}

	user, err := db.GetUser(i.UserId)
	if err != nil {
		return err
	}
	user.Name = iUserChange.Name
	user.Save()

	oUserChange := m_User{
		Id:   user.Id,
		Name: user.Name,
	}

	oPacket := mp(pt_OUserChange, oUserChange)
	Send(i.C, i.UserId, &oPacket)

	return nil
}
