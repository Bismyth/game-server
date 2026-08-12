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
func (p PlayerProperty[T]) MustGet(gameId uuid.UUID, playerId uuid.UUID) T {
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

func (p PlayerProperty[T]) GetMap(gameId uuid.UUID, players []uuid.UUID) (map[uuid.UUID]T, error) {
	o := make(map[uuid.UUID]T)
	for _, id := range players {
		v, err := p.Get(gameId, id)
		if err != nil {
			return o, err
		}
		o[id] = v
	}
	return o, nil
}

func MapApply[T any, U any](i map[uuid.UUID]T, cb func(v T) (U, error)) (map[uuid.UUID]U, error) {
	o := make(map[uuid.UUID]U)
	for id, ov := range i {
		v, err := cb(ov)
		if err != nil {
			return o, err
		}
		o[id] = v
	}
	return o, nil
}
