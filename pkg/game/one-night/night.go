package onenight

import (
	"encoding/json"
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
	role, err := GetPlayerProperty[Role](gameId, player, pd_initialRole)
	if err != nil {
		return err
	}
	if data.Role != role {
		return fmt.Errorf("submited action for unassigned role")
	}

	switch data.Role {
	case role_robber:
		err = handleRobber(gameId, player, data.Data)
	case role_seer:
		err = handleSeer(gameId, data.Data)
	}

	if err != nil {
		return err
	}

	return nil
}

func handleRobber(gameId uuid.UUID, player uuid.UUID, data json.RawMessage) error {
	var decoded RobberInput
	err := json.Unmarshal(data, &decoded)
	if err != nil {
		return err
	}

	decoded.Self = player

	err = SetProperty(gameId, d_robberInput, decoded)
	if err != nil {
		return err
	}

	return nil
}

func handleSeer(gameId uuid.UUID, data json.RawMessage) error {
	var decoded SeerInput
	err := json.Unmarshal(data, &decoded)
	if err != nil {
		return err
	}

	if decoded.SingleTarget == uuid.Nil && len(decoded.MultiTarget) != 2 {
		return fmt.Errorf("not target selected")
	}

	if decoded.SingleTarget == uuid.Nil {
		if decoded.MultiTarget[0] > 2 || decoded.MultiTarget[1] > 2 {
			return fmt.Errorf("selection out of range")
		}
	} else {
		err = targetInGame(gameId, decoded.SingleTarget)
		if err != nil {
			return err
		}
	}

	err = SetProperty(gameId, d_seerInput, decoded)
	if err != nil {
		return err
	}

	return nil
}

func targetInGame(gameId uuid.UUID, target uuid.UUID) error {
	inGame, err := db.IsUserInLobby(gameId, target)
	if err != nil {
		return err
	}
	if !inGame {
		return fmt.Errorf("target not in game")
	}

	return nil
}

//Night order
/*
- Seer
- Robber
- TroubleMaker
- Drunk


(info to send out)
- Werewolves
- Minion
- Masons
- Seer
- Robber
- Troublemaker
- Drunk
- Insomniac
*/

var executeFuncMap map[Role]func(gameId uuid.UUID) error = map[Role]func(gameId uuid.UUID) error{
	role_seer:   executeSeer,
	role_robber: executeRobber,
}

var executionOrder []Role = []Role{
	role_seer,
	role_robber,
	role_troublemaker,
	role_drunk,
}

func executeNight(c interfaces.GameCommunication, gameId uuid.UUID) func() {
	return func() {

		gameRoles, err := GetProperty[[]Role](gameId, d_roles)
		if err != nil {
			nightError(c, err)
			return
		}

		for _, role := range executionOrder {
			if roleInGame(gameRoles, role) {
				fn, ok := executeFuncMap[role]
				if !ok {
					nightError(c, fmt.Errorf("unrecognized role"))
				}
				err = fn(gameId)
				if err != nil {
					nightError(c, err)
					return
				}
			}
		}

		c.SendGlobal(&GameState{
			Public: &PublicGameState{
				NightOver: true,
			},
		})
	}
}

func executeSeer(gameId uuid.UUID) error {
	input, err := GetProperty[SeerInput](gameId, d_seerInput)
	if err != nil {
		return nil
	}

	var seerResult SeerData

	seerResult.Input = input

	if input.SingleTarget != uuid.Nil {
		role, err := truePlayerRole(gameId, input.SingleTarget)
		if err != nil {
			return err
		}
		seerResult.SingleResult = role
	} else if len(input.MultiTarget) != 2 {
		role1, err := GetProperty[Role](gameId, rolePos(input.MultiTarget[0]))
		if err != nil {
			return err
		}
		role2, err := GetProperty[Role](gameId, rolePos(input.MultiTarget[1]))
		if err != nil {
			return err
		}
		seerResult.MultiResult = []Role{role1, role2}
	}

	err = SetProperty(gameId, d_seerResult, seerResult)
	if err != nil {
		return err
	}

	return nil

}

func executeRobber(gameId uuid.UUID) error {
	input, err := GetProperty[RobberInput](gameId, d_robberInput)
	if err != nil {
		return nil
	}

	var output RobberData
	output.Input = input

	if input.Target != uuid.Nil && input.Self != uuid.Nil {

		targetRole, err := truePlayerRole(gameId, input.Target)
		if err != nil {
			return err
		}
		output.Stole = targetRole
		err = swapPlayers(gameId, input.Self, input.Target)
		if err != nil {
			return err
		}
	}
	err = SetProperty(gameId, d_robberResult, output)
	if err != nil {
		return err
	}

	return nil
}

func swapPlayers(gameId uuid.UUID, player1 uuid.UUID, player2 uuid.UUID) error {
	p1Pos, err := GetPlayerProperty[int](gameId, player1, pd_position)
	if err != nil {
		return err
	}

	p2Pos, err := GetPlayerProperty[int](gameId, player2, pd_position)
	if err != nil {
		return err
	}

	p1Role, err := GetProperty[Role](gameId, rolePos(p1Pos))
	if err != nil {
		return err
	}
	p2Role, err := GetProperty[Role](gameId, rolePos(p2Pos))
	if err != nil {
		return err
	}

	err = SetProperty(gameId, rolePos(p1Pos), p2Role)
	if err != nil {
		return err
	}

	err = SetProperty(gameId, rolePos(p2Pos), p1Role)
	if err != nil {
		return err
	}

	return nil
}
