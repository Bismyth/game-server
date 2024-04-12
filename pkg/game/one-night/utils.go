package onenight

import (
	"log"

	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

func handleNightError(c interfaces.GameCommunication, gameId uuid.UUID, err error) {
	// TODO: end game and broadcast what went wrong

	log.Println(err)
}
