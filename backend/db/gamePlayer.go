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

type Cursor struct {
	Key string
}

func (c Cursor) For(g uuid.UUID) BoundCursor {
	return BoundCursor{
		gameId: g,
		c:      c,
		conn:   getConn(),
		ctx:    context.Background(),
	}
}

type BoundCursor struct {
	gameId uuid.UUID
	c      Cursor
	conn   *redis.Client
	ctx    context.Context
}

func (bc BoundCursor) Key() string {
	return it(gameHashName, bc.gameId, bc.c.Key)
}

func (c BoundCursor) Add(id uuid.UUID) {
	err := c.conn.RPush(c.ctx, c.Key(), id.String()).Err()
	if err != nil {
		panic(err)
	}
}

func (c BoundCursor) GetAll() []uuid.UUID {
	idStrings, err := c.conn.LRange(c.ctx, c.Key(), 0, -1).Result()
	if err != nil {
		panic(err)
	}

	li, err := ParseUUIDList(idStrings)
	if err != nil {
		panic(err)
	}

	return li
}

func (c BoundCursor) Length() int64 {
	count, err := c.conn.LLen(c.ctx, c.Key()).Result()
	if err != nil {
		panic(err)
	}

	return count
}

func (c BoundCursor) RemoveTarget(playerId uuid.UUID) error {
	hasItem := c.HasItem(playerId)
	if !hasItem {
		return nil
	}
	current := c.Current()
	if current == playerId {
		c.Remove()
	} else {
		c.SeekIndex(playerId)
		c.Remove()
		c.SeekIndex(current)
	}

	return nil
}

func (c BoundCursor) HasItem(id uuid.UUID) bool {
	_, err := c.conn.LPos(c.ctx, c.Key(), id.String(), redis.LPosArgs{}).Result()

	return err == nil
}

func (c BoundCursor) Reset() {
	c.SetIndex(0)
}

func (c BoundCursor) GetIndex() int64 {
	s, err := c.conn.Get(c.ctx, ic(c.Key())).Result()
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

func (c BoundCursor) SetIndex(i int64) {
	err := c.conn.Set(c.ctx, ic(c.Key()), fmt.Sprintf("%d", i), 0).Err()
	if err != nil {
		log.Panic("failed to get cursor index")
	}
}

func (c BoundCursor) wrapIndex(i int64) int64 {
	size, err := c.conn.LLen(c.ctx, c.Key()).Result()
	if err != nil {
		log.Panic("failed to get size of index")
	}
	if size == 0 {
		return 0
	}

	return ((i % size) + size) % size
}

func (c BoundCursor) Next() uuid.UUID {
	c.Shift(1)
	return c.Current()
}

func (c BoundCursor) PeekNext() uuid.UUID {
	return c.PeekIndex(c.wrapIndex(c.GetIndex() + 1))
}

func (c BoundCursor) Previous() uuid.UUID {
	c.Shift(-1)
	return c.Current()
}

func (c BoundCursor) PeekPrevious() uuid.UUID {
	return c.PeekIndex(c.wrapIndex(c.GetIndex() - 1))
}

func (c BoundCursor) PeekIndex(i int64) uuid.UUID {
	idString, err := c.conn.LIndex(c.ctx, c.Key(), i).Result()
	if err != nil {
		log.Panic("could not get player index")
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		log.Panic("failed to parse player id")
	}

	return id
}

func (c BoundCursor) Current() uuid.UUID {
	return c.PeekIndex(c.GetIndex())
}

// Ends with the cursor on the next value
func (c BoundCursor) Remove() {
	v := c.Current()

	err := c.conn.LRem(c.ctx, c.Key(), 1, v.String()).Err()
	if err != nil {
		panic(err)
	}

	c.Shift(0)
}

func (c BoundCursor) SeekIndex(id uuid.UUID) {
	index, err := c.conn.LPos(c.ctx, c.Key(), id.String(), redis.LPosArgs{}).Result()
	if err != nil {
		panic(err)
	}

	c.SetIndex(index)
}

func (c BoundCursor) Shift(n int64) {
	c.SetIndex(c.wrapIndex(c.GetIndex() + n))
}

func (c BoundCursor) Delete() error {
	err := c.conn.Del(c.ctx, c.Key()).Err()
	if err != nil {
		return err
	}

	err = c.conn.Del(c.ctx, ic(c.Key())).Err()
	if err != nil {
		return err
	}

	return nil
}
