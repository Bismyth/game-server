package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func PlayerGiveType(gameId uuid.UUID, playerId uuid.UUID, playerType string) error {

	conn := getConn()
	ctx := context.Background()

	err := conn.RPush(ctx, it(gameHashName, gameId, playerType), playerId.String()).Err()
	if err != nil {
		return err
	}

	return nil
}

func PlayerRemType(gameId uuid.UUID, userId uuid.UUID, userType string) error {
	return fmt.Errorf("not implemented")
}

func PlayerTypeGetAll(gameId uuid.UUID, playerType string) ([]uuid.UUID, error) {
	conn := getConn()
	ctx := context.Background()

	idStrings, err := conn.LRange(ctx, it(gameHashName, gameId, playerType), 0, -1).Result()
	if err != nil {
		return nil, err
	}

	return ParseUUIDList(idStrings)
}

func PlayerTypeCount(gameId uuid.UUID, playerType string) (int64, error) {
	conn := getConn()
	ctx := context.Background()

	count, err := conn.LLen(ctx, it(gameHashName, gameId, playerType)).Result()
	if err != nil {
		return -1, err
	}

	return count, nil
}

func PlayerIsType(gameId uuid.UUID, playerId uuid.UUID, playerType string) bool {
	conn := getConn()
	ctx := context.Background()

	_, err := conn.LPos(ctx, it(gameHashName, gameId, playerType), playerId.String(), redis.LPosArgs{}).Result()

	return err == nil
}

func DeletePlayerTypeList(gameId uuid.UUID, playerType string) error {
	conn := getConn()
	ctx := context.Background()

	err := conn.Del(ctx, it(gameHashName, gameId, playerType)).Err()

	if err != nil {
		return err
	}

	return nil
}

type Cursor struct {
	key string
}

func GetCursor(gameId uuid.UUID, t string) *Cursor {
	c := &Cursor{
		key: it(gameHashName, gameId, t),
	}
	return c
}

func (c *Cursor) Reset() {
	c.SetIndex(0)
}

func (c *Cursor) GetIndex() int64 {
	conn := getConn()
	ctx := context.Background()

	s, err := conn.Get(ctx, ic(c.key)).Result()
	if errors.Is(err, redis.Nil) {
		return 0
	} else if err != nil {
		log.Panic("failed to get cursor index")
	}

	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		log.Panic("failed to get cursor index")
	}
	return i
}

func (c *Cursor) SetIndex(i int64) {
	conn := getConn()
	ctx := context.Background()

	err := conn.Set(ctx, ic(c.key), fmt.Sprintf("%d", i), 0).Err()
	if err != nil {
		log.Panic("failed to get cursor index")
	}
}

func (c *Cursor) wrapIndex(i int64) int64 {
	conn := getConn()
	ctx := context.Background()

	size, err := conn.LLen(ctx, c.key).Result()
	if err != nil {
		log.Panic("failed to get size of index")
	}
	if size == 0 {
		return 0
	}

	return ((i % size) + size) % size
}

func (c *Cursor) Next() (uuid.UUID, error) {
	c.Shift(1)
	return c.Current()
}

func (c *Cursor) PeekNext() (uuid.UUID, error) {
	return c.PeekIndex(c.wrapIndex(c.GetIndex() + 1))
}

func (c *Cursor) Previous() (uuid.UUID, error) {
	c.Shift(-1)
	return c.Current()
}

func (c *Cursor) PeekPrevious() (uuid.UUID, error) {
	return c.PeekIndex(c.wrapIndex(c.GetIndex() - 1))
}

func (c *Cursor) PeekIndex(i int64) (uuid.UUID, error) {
	conn := getConn()
	ctx := context.Background()

	idString, err := conn.LIndex(ctx, c.key, i).Result()
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not get player index")
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return id, fmt.Errorf("failed to parse player id")
	}

	return id, nil
}

func (c *Cursor) Current() (uuid.UUID, error) {
	return c.PeekIndex(c.GetIndex())
}

// Ends with the cursor on the next value
func (c *Cursor) Remove() error {
	conn := getConn()
	ctx := context.Background()

	v, err := c.Current()
	if err != nil {
		return err
	}

	err = conn.LRem(ctx, c.key, 1, v.String()).Err()
	if err != nil {
		return err
	}

	c.Shift(0)

	return nil
}

func (c *Cursor) SeekIndex(id uuid.UUID) error {
	conn := getConn()
	ctx := context.Background()

	index, err := conn.LPos(ctx, c.key, id.String(), redis.LPosArgs{}).Result()
	if err != nil {
		return err
	}

	c.SetIndex(index)

	return nil
}

func (c *Cursor) Shift(n int64) {
	c.SetIndex(c.wrapIndex(c.GetIndex() + n))
}

func (c *Cursor) Delete() error {
	conn := getConn()
	ctx := context.Background()

	err := conn.Del(ctx, c.key).Err()
	if err != nil {
		return err
	}

	err = conn.Del(ctx, ic(c.key)).Err()
	if err != nil {
		return err
	}

	return nil
}
