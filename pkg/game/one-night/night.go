package onenight

import (
	"fmt"
	"time"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

func startNight(c interfaces.GameCommunication, gameId uuid.UUID, player uuid.UUID) error {
	isHost, err := db.IsUserLobbyHost(gameId, player)
	if err != nil {
		return err
	}
	if !isHost {
		return fmt.Errorf("not authorised to start game")
	}

	nightDuration, err := GetProperty[int](gameId, d_nightTime)
	if err != nil {
		return err
	}

	time.AfterFunc(time.Second*time.Duration(nightDuration), executeNight(c, gameId))

	return nil
}

func handleRoleAction(gameId uuid.UUID, player uuid.UUID, data RoleAction) error {
	return nil
}

func executeNight(c interfaces.GameCommunication, gameId uuid.UUID) func() {
	return func() {
		c.SendGlobal(&GameState{
			Public: &PublicGameState{
				NightOver: true,
			},
		})
	}
}
