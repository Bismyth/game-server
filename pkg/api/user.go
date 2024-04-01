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
	name, err := db.GetUserName(userId)
	if err != nil {
		log.Panic("Failed to initilize user")
	}

	oData := m_User{
		Id:   userId,
		Name: name,
	}

	oPacket := mp(pt_OUserInit, oData)

	return MarshalPacket(&oPacket)
}

// User change event
const pt_IUserChange IPacketType = "client_user_change"
const pt_OUserChange OPacketType = "server_user_change"

func handleUserChange(i HandlerInput) error {
	iUserChange, err := hp[m_User](i.Packet)
	if err != nil {
		return err
	}

	err = db.SetUserName(i.UserId, iUserChange.Name)
	if err != nil {
		return err
	}

	oUserChange := m_User{
		Id:   i.UserId,
		Name: iUserChange.Name,
	}

	oPacket := mp(pt_OUserChange, oUserChange)
	Send(i.C, i.UserId, &oPacket)

	return nil
}
