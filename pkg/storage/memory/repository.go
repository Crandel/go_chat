package memory

import (
	"sync"
)

type Storage struct {
	Users    map[UserId]User
	Rooms    map[string]Room
	Messages map[int]Message
	*sync.RWMutex
}

func NewStorage() Storage {
	return Storage{
		Users:    map[UserId]User{},
		Rooms:    map[string]Room{},
		Messages: map[int]Message{},
		RWMutex:  &sync.RWMutex{},
	}
}
