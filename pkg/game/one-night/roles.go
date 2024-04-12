package onenight

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
