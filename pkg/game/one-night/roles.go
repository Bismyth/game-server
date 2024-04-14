package onenight

import "github.com/google/uuid"

type Role string

const role_werewolf = "werewolf"
const role_villager = "villager"
const role_robber = "robber"
const role_drunk = "drunk"
const role_troublemaker = "troublemaker"
const role_mason = "mason"
const role_minion = "minion"
const role_insominiac = "insomniac"
const role_seer = "seer"

var allRoles = []Role{
	role_werewolf,
	role_villager,
	role_robber,
	role_drunk,
	role_troublemaker,
	role_mason,
	role_minion,
	role_insominiac,
	role_seer,
}

func isValidRole(r Role) bool {
	for _, a := range allRoles {
		if a == r {
			return true
		}
	}
	return false
}

type RoleInput[T any] struct {
	Self uuid.UUID
}

type RobberInput struct {
	Target uuid.UUID `json:"target"`
	Self   uuid.UUID `json:"self"`
}

type SeerInput struct {
	SingleTarget uuid.UUID `json:"singleTarget"`
	MultiTarget  []int     `json:"multiTarget"`
}

type DrunkInput struct {
	Target int
}

type TroubleMakerInput struct {
	Target1 uuid.UUID
	Target2 uuid.UUID
}

type RobberData struct {
	Input RobberInput `json:"input"`
	Stole Role        `json:"stole"`
}

type SeerData struct {
	Input        SeerInput `json:"input"`
	SingleResult Role      `json:"singleResult"`
	MultiResult  []Role    `json:"multiResult"`
}

type WerewolfData struct {
	Wolves []uuid.UUID
}

type MasonData struct {
	Masons []uuid.UUID
}

type InsomniacData struct {
	Result Role
}
