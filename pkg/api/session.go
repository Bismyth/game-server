package api

import (
	"fmt"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
)

func HandleSessionInit(c Client, claims RoomTokenClaims, sessionId uuid.UUID) error {
	exists, err := db.RoomHasUser(claims.RoomId, claims.UserId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("not in room")
	}

	err = db.SetRoomUserSession(claims.RoomId, claims.UserId, sessionId)
	if err != nil {
		return err
	}

	err = db.NewSession(sessionId, claims.RoomId, claims.UserId)
	if err != nil {
		return err
	}

	sendRoomInfoSingle(c, claims.RoomId, sessionId)

	return nil
}

func HandleSessionClose(sessionId uuid.UUID) error {
	d, err := db.GetSessionDetails(sessionId)
	if err != nil {
		return err
	}
	err = db.RemoveRoomUserSession(d.RoomId, d.UserId)
	if err != nil {
		return err
	}

	err = db.DeleteSession(sessionId)
	if err != nil {
		return err
	}

	return nil
}
