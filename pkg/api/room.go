package api

import (
	"fmt"
	"log"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/game"
	"github.com/google/uuid"
)

type m_Room struct {
	MaxPlayers *int        `json:"maxPlayers"`
	Name       *string     `json:"name"`
	Host       *uuid.UUID  `json:"host"`
	GameType   *string     `json:"gameType"`
	InGame     *bool       `json:"inGame"`
	Options    interface{} `json:"gameOptions"`
}

type m_RoomChange struct {
	MaxPlayers int
	GameType   string
}

const pt_ORoomUserChange = "server_room_user_change"
const pt_ORoomInfo = "server_room_info"

type m_RoomUserList struct {
	Order []uuid.UUID `json:"order"`
	Names db.UserList `json:"names"`
}

func makeRoomUserMessage(roomId uuid.UUID) *Packet[m_RoomUserList] {
	users, err := db.GetRoomUsers(roomId)
	if err != nil {
		return nil
	}
	var userList m_RoomUserList
	userList.Names = users

	order, err := db.GetRoomUserOrder(roomId)
	if err != nil {
		return nil
	}
	userList.Order = order

	newMessage := mp(pt_ORoomUserChange, userList)
	return &newMessage
}

func sendRoomUsers(c Client, roomId uuid.UUID) {
	sessions, err := db.GetRoomSessions(roomId)
	if err != nil {
		log.Printf("failed to get sessions for room: %v", err)
		return
	}
	packet := makeRoomUserMessage(roomId)
	SendMany(c, sessions, packet)
}

func roomInfoPacket(roomId uuid.UUID) (*Packet[m_Room], error) {
	properties := []string{"name", "maxPlayers", "gameType", "inGame", "options", "host"}
	m, err := db.GetRoomProperties(roomId, properties)
	if err != nil {
		return nil, err
	}

	var info m_Room
	err = decode(m, &info)
	if err != nil {
		return nil, err
	}

	packet := mp(pt_ORoomInfo, info)

	return &packet, nil
}

func sendRoomInfoSingle(c Client, roomId uuid.UUID, sessionId uuid.UUID) {
	roomInfo, err := roomInfoPacket(roomId)
	if err != nil {
		log.Println("failed to make lobby info packet")
	}

	Send(c, sessionId, roomInfo)
}

func sendRoomInfo(c Client, roomId uuid.UUID) {
	roomInfo, err := roomInfoPacket(roomId)
	if err != nil {
		log.Println("failed to make lobby info packet")
	}

	ids, err := db.GetRoomSessions(roomId)
	if err != nil {
		log.Println("failed to get room sessions")
	}

	SendMany(c, ids, roomInfo)
}

func makeChangesAllowed(roomId uuid.UUID, userId uuid.UUID) error {
	hostId, err := db.GetRoomProperty[uuid.UUID](roomId, "host")
	if err != nil {
		return err
	}
	if hostId != userId {
		return fmt.Errorf("not room host")
	}

	return nil
}

const pt_IRoomChange = "client_room_change"

func roomChange(i HandlerInput) error {
	data, err := hp[m_RoomChange](i.Packet)
	if err != nil {
		return err
	}

	err = makeChangesAllowed(i.Session.RoomId, i.Session.UserId)
	if err != nil {
		return err
	}

	oldGameType, err := db.GetRoomProperty[string](i.Session.RoomId, "gameType")
	if err != nil {
		return err
	}
	if data.GameType != oldGameType {
		options, err := game.GetDefaultOptions(data.GameType)
		if err != nil {
			return err
		}

		err = db.SetRoomProperty(i.Session.RoomId, "gameType", data.GameType)
		if err != nil {
			return err
		}

		err = db.SetRoomProperty(i.Session.RoomId, "options", options)
		if err != nil {
			return err
		}
	}

	sendRoomInfo(i.C, i.Session.RoomId)

	return nil
}

const pt_IRoomUsers = "client_room_users"

func roomUsers(i HandlerInput) error {
	packet := makeRoomUserMessage(i.Session.RoomId)
	Send(i.C, i.SessionId, packet)

	return nil
}

const pt_IRoomOptions = "client_room_options"

func roomOptions(i HandlerInput) error {
	err := makeChangesAllowed(i.Session.RoomId, i.Session.UserId)
	if err != nil {
		return err
	}

	err = db.SetRoomProperty(i.Session.RoomId, "options", string(i.Packet.Data))
	if err != nil {
		return err
	}

	sendRoomInfo(i.C, i.Session.RoomId)
	return nil
}

const pt_IRoomLeave = "client_room_leave"

func roomLeave(i HandlerInput) error {
	err := handleLeave(i.C, i.Session.RoomId, i.Session.UserId)
	if err != nil {
		return err
	}
	i.C.Close(i.SessionId)
	return nil
}

func handleLeave(c Client, roomId uuid.UUID, userId uuid.UUID) error {
	ig, err := db.GetRoomProperty[bool](roomId, "inGame")
	if err != nil {
		return err
	}
	if ig {
		g := &gc{
			RoomId: roomId,
			C:      c,
		}
		err = game.HandleLeave(g, roomId, userId)
		if err != nil {
			return err
		}
	}

	err = db.RemoveRoomUser(roomId, userId)
	if err != nil {
		return err
	}

	sendRoomInfo(c, roomId)
	sendRoomUsers(c, roomId)

	return nil
}

const pt_IRoomKick = "client_room_kick"
const pt_ORoomKick = "server_room_kick"

func roomKick(i HandlerInput) error {
	err := makeChangesAllowed(i.Session.RoomId, i.Session.UserId)
	if err != nil {
		return err
	}

	id, err := hp[uuid.UUID](i.Packet)
	if err != nil {
		return err
	}

	session, err := db.GetRoomUserSession(i.Session.RoomId, *id)
	if err != nil {
		return err
	}

	packet := mp(pt_ORoomKick, "")
	Send(i.C, session, &packet)

	i.C.Close(session)

	err = handleLeave(i.C, i.Session.RoomId, *id)
	if err != nil {
		return err
	}

	return nil
}
