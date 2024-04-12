package onenight

import (
	"fmt"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
)

type Options struct {
	Roles     []Role `json:"roles"`
	NightTime int    `json:nightTime`
}

type GameState struct {
	Public  *PublicGameState  `json:"public"`
	Private *PrivateGameState `json:"private"`
}

type PublicGameState struct {
}

type PrivateGameState struct {
	Role     Role   `json:"role"`
	RoleInfo []byte `json:"roleInfo"`
}

type DBProperty string

const d_nightTime = "night"

type DBPlayerProprety string

const pd_position DBPlayerProprety = "position"
const pd_data DBPlayerProprety = "data"

func rolePos(i int) DBProperty {
	return DBProperty(fmt.Sprintf("role:%d", i))
}

func GetProperty[T any](gameId uuid.UUID, p DBProperty) (T, error) {
	return db.GetGameProperty[T](gameId, string(p))
}

func SetProperty[T any](gameId uuid.UUID, p DBProperty, data T) error {
	return db.SetGameProperty(gameId, string(p), data)
}

func SetPlayerProperty[T any](gameId uuid.UUID, playerId uuid.UUID, p DBPlayerProprety, data T) error {
	return db.SetPlayerProperty(gameId, playerId, string(p), data)
}

func GetPlayerProperty[T any](gameId uuid.UUID, playerId uuid.UUID, p DBPlayerProprety) (T, error) {
	return db.GetPlayerProperty[T](gameId, playerId, string(p))
}
