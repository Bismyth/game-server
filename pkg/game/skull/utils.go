package skull

import (
	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
)

func cleanup(gameId uuid.UUID) error {
	c := db.GetCursor(gameId, playerType)
	err := c.Delete()
	if err != nil {
		return err
	}

	err = db.ExpireCache(gameId, cacheExpireTime)
	if err != nil {
		return err
	}

	err = db.DeletePlayerTypeList(gameId, playerType)
	if err != nil {
		return err
	}

	err = db.DeleteGame(gameId)
	if err != nil {
		return err
	}

	return nil
}
