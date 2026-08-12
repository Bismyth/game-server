package db

import (
	"context"
	"math/rand/v2"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Deck[T any] struct {
	Key string
}

func (d Deck[T]) For(gId uuid.UUID) BoundDeck[T] {
	return BoundDeck[T]{
		keyPart: d.Key,
		gameId:  gId,
		conn:    getConn(),
		ctx:     context.Background(),
	}
}

type BoundDeck[T any] struct {
	keyPart string
	gameId  uuid.UUID
	conn    *redis.Client
	ctx     context.Context
}

func (d BoundDeck[T]) Key() string {
	return it(gameHashName, d.gameId, d.keyPart)
}

func (d BoundDeck[T]) Add(v T) {
	value, err := Encode(v)
	if err != nil {
		panic(err)
	}

	err = d.conn.LPush(d.ctx, d.Key(), value).Err()
	if err != nil {
		panic(err)
	}
}

func (d BoundDeck[T]) Draw() T {
	raw, err := d.conn.LPop(d.ctx, d.Key()).Bytes()
	if err != nil {
		panic(err)
	}

	var value T
	err = Decode(raw, &value)
	if err != nil {
		panic(err)
	}
	return value
}

func (d BoundDeck[T]) Shuffle() {
	a := int(d.Length())
	arrCopy := make([]T, a)
	for i := 0; i < a; i++ {
		arrCopy[i] = d.Draw()
	}

	for i := a - 1; i > 0; i-- {
		j := rand.IntN(i + 1)
		arrCopy[i], arrCopy[j] = arrCopy[j], arrCopy[i]
	}

	for _, v := range arrCopy {
		d.Add(v)
	}
}

func (d BoundDeck[T]) Length() int64 {
	length, err := d.conn.LLen(d.ctx, d.Key()).Result()
	if err != nil {
		panic(err)
	}

	return length
}

func (d BoundDeck[T]) GetAll() []T {
	strings, err := d.conn.LRange(d.ctx, d.Key(), 0, -1).Result()
	if err != nil {
		panic(err)
	}

	a := len(strings)
	output := make([]T, a)
	for i, s := range strings {
		var v T
		err := Decode([]byte(s), &v)
		if err != nil {
			panic(err)
		}
		output[i] = v
	}

	return output
}

func (d BoundDeck[T]) Clear() {
	err := d.conn.Del(d.ctx, d.Key()).Err()
	if err != nil {
		panic(err)
	}
}
