package onenight

import (
	"github.com/google/uuid"
)

func getPrivateInfo(gameId, playerId uuid.UUID) (*PrivateGameState, error) {
	var pi PrivateGameState

	pos, err := GetPlayerProperty[int](gameId, playerId, pd_position)
	if err != nil {
		return nil, err
	}

	role, err := GetProperty[Role](gameId, rolePos(pos))
	if err != nil {
		return nil, err
	}
	roleData, err := GetPlayerProperty[string](gameId, playerId, pd_data)
	if err != nil {
		return nil, err
	}

	pi.Role = role
	pi.RoleInfo = []byte(roleData)

	return &pi, nil

}
