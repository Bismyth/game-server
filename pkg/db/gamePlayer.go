package db

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
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
	return nil
}

func PlayerTypeGetAll(gameId uuid.UUID, playerType string) ([]uuid.UUID, error) {
	return nil, nil
}

type Cursor struct {
	index int64
	key   string
}

func NewCursor(gameId uuid.UUID, t string) *Cursor {
	return &Cursor{
		index: 0,
		key:   it(gameHashName, gameId, t),
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
	return c.PeekIndex(c.wrapIndex(c.index + 1))
}

func (c *Cursor) Previous() (uuid.UUID, error) {
	c.Shift(-1)
	return c.Current()
}

func (c *Cursor) PeekPrevious() (uuid.UUID, error) {
	return c.PeekIndex(c.wrapIndex(c.index - 1))
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
	return c.PeekIndex(c.index)
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

func (c *Cursor) Shift(n int64) {
	c.index = c.wrapIndex(c.index + n)
}
