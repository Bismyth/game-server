package db

import (
	"context"
	"log"

	"github.com/Bismyth/game-server/db/msg"
	"github.com/Bismyth/game-server/interfaces"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const eventHashName = "events"

func GameEvent(gameId uuid.UUID, c interfaces.GameCommunication) BoundEvent {
	return BoundEvent{gameId: gameId, c: c}
}

type BoundEvent struct {
	gameId uuid.UUID
	c      interfaces.GameCommunication
}

func (be BoundEvent) Log(data msg.EventData) {
	conn := getConn()
	ctx := context.Background()

	_, err := conn.XAdd(ctx, &redis.XAddArgs{
		Stream: i(eventHashName, be.gameId),
		Values: map[string]string{"message": string(data)},
	}).Result()
	if err != nil {
		log.Printf("failed to write event to redis: %v", err)
	}

	be.c.SendEvent(data)
}
