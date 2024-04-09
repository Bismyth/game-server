package api

import (
	"fmt"
	"log"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/game"
	"github.com/google/uuid"
)

const pt_OGameEvent OPacketType = "server_game_event"
const pt_OGameState OPacketType = "server_game_state"
const pt_OGameAction OPacketType = "server_game_action"

const pt_IGameAction IPacketType = "client_game_action"
const pt_IGameNew IPacketType = "client_game_new"
const pt_IGameReady IPacketType = "client_game_ready"

type gc struct {
	GameId uuid.UUID
	C      Client
}

func (g *gc) SendEvent(data any) {
	players, err := db.GetLobbyUserIds(g.GameId)
	if err != nil {
		log.Println("could not send event as players couldnt be retrieved")
	}

	packet := mp(pt_OGameEvent, data)

	SendMany(g.C, players, &packet)
}

func (g *gc) SendGlobal(data any) {
	players, err := db.GetLobbyUserIds(g.GameId)
	if err != nil {
		log.Println("could not send event as players couldnt be retrieved")
	}

	packet := mp(pt_OGameState, data)

	SendMany(g.C, players, &packet)
}

func (g *gc) SendPlayer(playerId uuid.UUID, data any) {
	packet := mp(pt_OGameState, data)
	Send(g.C, playerId, &packet)
}

func (g *gc) ActionPrompt(playerId uuid.UUID, data any) {
	packet := mp(pt_OGameAction, data)
	Send(g.C, playerId, &packet)
}

func (g *gc) EndGame() {
	err := db.SetLobbyProperty(g.GameId, "inGame", false)
	if err != nil {
		log.Printf("failed to change lobby to not in game")
	}
	sendLobbyData(g.C, g.GameId)
}

func gameNew(i HandlerInput) error {
	lobbyId, err := hp[uuid.UUID](i.Packet)
	if err != nil {
		return err
	}

	err = game.New(*lobbyId)
	if err != nil {
		SendGameErr(i.C, i.UserId, err)
		return nil
	}

	sendLobbyData(i.C, *lobbyId)

	return nil
}

type SharedGameData struct {
	Id uuid.UUID
}

func gameReady(i HandlerInput) error {
	data, err := hp[SharedGameData](i.Packet)
	if err != nil {
		SendGameErr(i.C, i.UserId, fmt.Errorf("failed to unmarshal input"))
	}

	g := &gc{
		GameId: data.Id,
		C:      i.C,
	}

	err = game.HandleReady(g, i.UserId, i.Packet.Data)
	if err != nil {
		SendGameErr(i.C, i.UserId, err)
	}

	return nil
}

func gameAction(i HandlerInput) error {
	data, err := hp[SharedGameData](i.Packet)
	if err != nil {
		SendGameErr(i.C, i.UserId, fmt.Errorf("failed to unmarshal input"))
	}

	g := &gc{
		GameId: data.Id,
		C:      i.C,
	}

	err = game.HandleAction(g, i.UserId, i.Packet.Data)
	if err != nil {
		SendGameErr(i.C, i.UserId, err)
	}

	return nil
}
