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

func FilledStorage(users map[UserId]User, rooms map[string]Room, messages map[int]Message) Storage {
	return Storage{
		Users:    users,
		Rooms:    rooms,
		Messages: messages,
		RWMutex:  &sync.RWMutex{},
	}
}
