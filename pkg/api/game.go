package api

import (
	"log"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/game"
	"github.com/google/uuid"
)

const pt_OGameEvent OPacketType = "server_game_event"
const pt_OGameState OPacketType = "server_game_state"
const pt_OGameAction OPacketType = "server_game_action"

const pt_IGameAction IPacketType = "client_game_action"
const pt_IGameStart IPacketType = "client_game_start"
const pt_IGameReady IPacketType = "client_game_ready"

type gc struct {
	RoomId uuid.UUID
	C      Client
}

func (g *gc) SendEvent(data any) {
	sessions, err := db.GetRoomSessions(g.RoomId)
	if err != nil {
		log.Println("could not send event as players couldnt be retrieved")
	}

	packet := mp(pt_OGameEvent, data)

	SendMany(g.C, sessions, &packet)
}

func (g *gc) SendGlobal(data any) {
	sessions, err := db.GetRoomSessions(g.RoomId)
	if err != nil {
		log.Println("could not send event as players couldnt be retrieved")
	}

	packet := mp(pt_OGameState, data)

	SendMany(g.C, sessions, &packet)
}

func (g *gc) SendPlayer(playerId uuid.UUID, data any) {
	session, err := db.GetRoomUserSession(g.RoomId, playerId)
	if err != nil {
		return
	}

	packet := mp(pt_OGameState, data)
	Send(g.C, session, &packet)
}

func (g *gc) ActionPrompt(playerId uuid.UUID, data any) {
	session, err := db.GetRoomUserSession(g.RoomId, playerId)
	if err != nil {
		return
	}

	packet := mp(pt_OGameAction, data)
	Send(g.C, session, &packet)
}

func (g *gc) EndGame() {
	err := db.SetRoomProperty(g.RoomId, "inGame", false)
	if err != nil {
		log.Printf("failed to change lobby to not in game")
	}
	sendRoomInfo(g.C, g.RoomId)
}

func gameStart(i HandlerInput) error {
	err := makeChangesAllowed(i.Session.RoomId, i.Session.UserId)
	if err != nil {
		return err
	}

	err = game.New(i.Session.RoomId)
	if err != nil {
		SendGameErr(i.C, i.SessionId, err)
		return nil
	}

	sendRoomInfo(i.C, i.Session.RoomId)

	return nil
}

func gameReady(i HandlerInput) error {
	g := &gc{
		RoomId: i.Session.RoomId,
		C:      i.C,
	}

	err := game.HandleReady(g, i.Session.RoomId, i.Session.UserId, i.Packet.Data)
	if err != nil {
		SendGameErr(i.C, i.SessionId, err)
	}

	return nil
}

func gameAction(i HandlerInput) error {
	g := &gc{
		RoomId: i.Session.RoomId,
		C:      i.C,
	}

	err := game.HandleAction(g, i.Session.RoomId, i.Session.UserId, i.Packet.Data)
	if err != nil {
		SendGameErr(i.C, i.SessionId, err)
	}

	return nil
}
