package db

import "github.com/google/uuid"

type Property[T any] struct {
	Key string
}

func (p Property[T]) Get(gameId uuid.UUID) (T, error) {
	return GetGameProperty[T](gameId, p.Key)
}
func (p Property[T]) MustGet(gameId uuid.UUID) T {
	r, err := GetGameProperty[T](gameId, p.Key)
	if err != nil {
		panic(err)
	}
	return r
}

func (p Property[T]) Set(gameId uuid.UUID, data T) error {
	return SetGameProperty(gameId, p.Key, data)
}

func (p Property[T]) MustSet(gameId uuid.UUID, data T) {
	err := SetGameProperty(gameId, p.Key, data)
	if err != nil {
		panic(err)
	}
}
