package onenight

import (
	"github.com/google/uuid"
)

func getPrivateInfo(gameId, playerId uuid.UUID) (*PrivateGameState, error) {
	var pi PrivateGameState

	role, err := GetPlayerProperty[Role](gameId, playerId, pd_initialRole)
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
