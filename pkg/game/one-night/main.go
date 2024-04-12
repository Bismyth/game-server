package onenight

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

const Code = "onenightwerewolf"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) DefaultOptions() interface{} {
	return &Options{
		Roles:     []Role{role_werewolf, role_werewolf, role_robber, role_seer, role_villager, role_villager},
		NightTime: 60,
	}
}

func (h *Handler) New(gameId uuid.UUID, options []byte) error {
	var parsedOptions Options
	err := json.Unmarshal(options, &parsedOptions)
	if err != nil {
		return err
	}

	if parsedOptions.NightTime <= 0 {
		return fmt.Errorf("night time not set")
	}

	if !ValidateRoleAmounts(parsedOptions.Roles) {
		return fmt.Errorf("invalid role set given")
	}

	players, err := db.GetLobbyUserIds(gameId)
	if err != nil {
		return err
	}

	if len(players) < 3 {
		return fmt.Errorf("not enough players to start game")
	}

	if len(players)+3 != len(parsedOptions.Roles) {
		return fmt.Errorf("for this game you would need %d people you have: %d", len(parsedOptions.Roles)-3, len(players))
	}

	//shuffle roles
	rand.Shuffle(len(parsedOptions.Roles), func(i, j int) {
		parsedOptions.Roles[i], parsedOptions.Roles[j] = parsedOptions.Roles[j], parsedOptions.Roles[i]
	})

	for i, role := range parsedOptions.Roles {
		if !isValidRole(role) {
			return fmt.Errorf("unrecognized role included")
		}
		err = SetProperty(gameId, rolePos(i), role)
		if err != nil {
			return err
		}
	}

	for i, player := range players {
		err = SetPlayerProperty(gameId, player, pd_position, i+3)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateRoleAmounts(roles []Role) bool {
	//TODO: validate whether any role has too many cards
	return true
}

func (h *Handler) HandleReady(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error {
	pi, err := getPrivateInfo(gameId, playerId)
	if err != nil {
		return err
	}

	c.SendPlayer(playerId, &GameState{
		Private: pi,
	})

	return nil
}

func (h *Handler) HandleAction(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID, data json.RawMessage) error {

	return fmt.Errorf("not implemented")
}

func (h *Handler) HandleLeave(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error {
	return fmt.Errorf("not implemented")
}
