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
	return BoundCursor{gameId: g, c: c}
}

type BoundCursor struct {
	gameId uuid.UUID
	c      Cursor
}

func (bc BoundCursor) Key() string {
	return it(gameHashName, bc.gameId, bc.c.Key)
}

func (c BoundCursor) Add(id uuid.UUID) error {
	conn := getConn()
	ctx := context.Background()

	err := conn.RPush(ctx, c.Key(), id.String()).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c BoundCursor) GetAll() ([]uuid.UUID, error) {
	conn := getConn()
	ctx := context.Background()

	idStrings, err := conn.LRange(ctx, c.Key(), 0, -1).Result()
	if err != nil {
		return nil, err
	}

	return ParseUUIDList(idStrings)
}

func (c BoundCursor) Length() (int64, error) {
	conn := getConn()
	ctx := context.Background()

	count, err := conn.LLen(ctx, c.Key()).Result()
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (c BoundCursor) RemoveTarget(playerId uuid.UUID) error {
	hasItem := c.HasItem(playerId)
	if !hasItem {
		return nil
	}
	current, err := c.Current()
	if err != nil {
		return err
	}
	if current == playerId {
		err := c.Remove()
		if err != nil {
			return err
		}
	} else {
		err := c.SeekIndex(playerId)
		if err != nil {
			return err
		}
		err = c.Remove()
		if err != nil {
			return err
		}
		err = c.SeekIndex(current)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c BoundCursor) HasItem(id uuid.UUID) bool {
	conn := getConn()
	ctx := context.Background()

	_, err := conn.LPos(ctx, c.Key(), id.String(), redis.LPosArgs{}).Result()

	return err == nil
}

func (c BoundCursor) Reset() {
	c.SetIndex(0)
}

func (c BoundCursor) GetIndex() int64 {
	conn := getConn()
	ctx := context.Background()

	s, err := conn.Get(ctx, ic(c.Key())).Result()
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
	conn := getConn()
	ctx := context.Background()

	err := conn.Set(ctx, ic(c.Key()), fmt.Sprintf("%d", i), 0).Err()
	if err != nil {
		log.Panic("failed to get cursor index")
	}
}

func (c BoundCursor) wrapIndex(i int64) int64 {
	conn := getConn()
	ctx := context.Background()

	size, err := conn.LLen(ctx, c.Key()).Result()
	if err != nil {
		log.Panic("failed to get size of index")
	}
	if size == 0 {
		return 0
	}

	return ((i % size) + size) % size
}

func (c BoundCursor) Next() (uuid.UUID, error) {
	c.Shift(1)
	return c.Current()
}

func (c BoundCursor) PeekNext() (uuid.UUID, error) {
	return c.PeekIndex(c.wrapIndex(c.GetIndex() + 1))
}

func (c BoundCursor) Previous() (uuid.UUID, error) {
	c.Shift(-1)
	return c.Current()
}

func (c BoundCursor) PeekPrevious() (uuid.UUID, error) {
	return c.PeekIndex(c.wrapIndex(c.GetIndex() - 1))
}

func (c BoundCursor) PeekIndex(i int64) (uuid.UUID, error) {
	conn := getConn()
	ctx := context.Background()

	idString, err := conn.LIndex(ctx, c.Key(), i).Result()
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not get player index")
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return id, fmt.Errorf("failed to parse player id")
	}

	return id, nil
}

func (c BoundCursor) Current() (uuid.UUID, error) {
	return c.PeekIndex(c.GetIndex())
}

// Ends with the cursor on the next value
func (c BoundCursor) Remove() error {
	conn := getConn()
	ctx := context.Background()

	v, err := c.Current()
	if err != nil {
		return err
	}

	err = conn.LRem(ctx, c.Key(), 1, v.String()).Err()
	if err != nil {
		return err
	}

	c.Shift(0)

	return nil
}

func (c BoundCursor) SeekIndex(id uuid.UUID) error {
	conn := getConn()
	ctx := context.Background()

	index, err := conn.LPos(ctx, c.Key(), id.String(), redis.LPosArgs{}).Result()
	if err != nil {
		return err
	}

	c.SetIndex(index)

	return nil
}

func (c BoundCursor) Shift(n int64) {
	c.SetIndex(c.wrapIndex(c.GetIndex() + n))
}

func (c BoundCursor) Delete() error {
	conn := getConn()
	ctx := context.Background()

	err := conn.Del(ctx, c.Key()).Err()
	if err != nil {
		return err
	}

	err = conn.Del(ctx, ic(c.Key())).Err()
	if err != nil {
		return err
	}

	return nil
}
