package db

import (
	"github.com/google/uuid"
)

type PlayerProperty[T any] struct {
	Key string
}

func (p PlayerProperty[T]) Get(gameId uuid.UUID, playerId uuid.UUID) (T, error) {
	return GetPlayerProperty[T](gameId, playerId, p.Key)
}
func (p PlayerProperty[T]) MustGetPD(gameId uuid.UUID, playerId uuid.UUID) T {
	d, err := p.Get(gameId, playerId)
	if err != nil {
		panic(err)
	}
	return d
}

func (p PlayerProperty[T]) Set(gameId uuid.UUID, playerId uuid.UUID, data T) error {
	return SetPlayerProperty(gameId, playerId, p.Key, data)
}
func (p PlayerProperty[T]) MustSet(gameId uuid.UUID, playerId uuid.UUID, data T) {
	err := SetPlayerProperty(gameId, playerId, p.Key, data)
	if err != nil {
		panic(err)
	}
}

func (p PlayerProperty[T]) GetMulti(gameId uuid.UUID, players []uuid.UUID) ([]T, error) {
	o := make([]T, len(players))
	for i, id := range players {
		v, err := p.Get(gameId, id)
		if err != nil {
			return o, err
		}
		o[i] = v
	}
	return o, nil
}

func LoadAllProperty[T any, U any](gameId uuid.UUID, players []uuid.UUID, p PlayerProperty[T], cb func(T) (U, error)) (map[uuid.UUID]U, error) {
	return LoadPlayerProperty(gameId, players, p.Key, cb)
}
