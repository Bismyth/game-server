package onenight

import (
	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

func roleInGame(allRoles []Role, role Role) bool {
	for _, a := range allRoles {
		if a == role {
			return true
		}
	}
	return false
}

func truePlayerRole(gameId uuid.UUID, target uuid.UUID) (Role, error) {
	pos, err := GetPlayerProperty[int](gameId, target, pd_position)
	if err != nil {
		return "", err
	}
	role, err := GetProperty[Role](gameId, rolePos(pos))
	if err != nil {
		return "", err
	}

	return role, nil
}

func nightError(c interfaces.GameCommunication, err error) {
	// TODO: do something to handle night execution error
}
