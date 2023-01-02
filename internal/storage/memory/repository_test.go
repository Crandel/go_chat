package memory_test

import (
	"sync"

	m "github.com/Crandel/go_chat/internal/storage/memory"
)

func NewTestStorage() *m.Storage {
	return &m.Storage{
		Users:    make(map[m.UserId]m.User, 0),
		Rooms:    make(map[string]m.Room, 0),
		Messages: make(map[int]m.Message, 0),
		RWMutex:  &sync.RWMutex{},
	}
}
