package memory

import (
	"math/rand"
	"time"

	cht "github.com/Crandel/go_chat/internal/chatting"
	lg "github.com/Crandel/go_chat/internal/logging"
)

func (str *Storage) WriteMessage(u *cht.Client, r *cht.Room, msg string) error {
	rand.Seed(time.Now().UnixNano())
	id := rand.Int()
	str.Messages[id] = Message{
		ID:       id,
		UserID:   UserId(u.Nick),
		RoomName: r.Name,
		Payload:  msg,
		Created:  time.Now(),
	}
	return nil
}

func (str *Storage) ExcludeFromRoom(roomName string, u *cht.Client) error {
	const op lg.Op = "memory.ExcludeFromRoom"
	var mr Room
	mr, exists := str.Rooms[roomName]
	if !exists {
		return lg.New(op, "No room with name "+roomName)
	}
	for i, ru := range mr.Members {
		if string(ru) == u.Nick {
			mr.Members = append(mr.Members[:i], mr.Members[i+1:]...)
		}
	}
	return nil
}

func (str *Storage) AddUserToRoom(name string, c *cht.Client) error {
	const op lg.Op = "memory.AddUserToRoom"
	var mr Room
	mr, exists := str.Rooms[name]
	if !exists {
		mrname, error := str.AddRoom(name)
		if error != nil {
			return error
		}
		mr = str.Rooms[mrname]
	}
	ru, exists := str.Users[UserId(c.Nick)]
	if !exists {
		return lg.New(op, "No user with id "+c.Nick)
	}
	mr.Members = append(mr.Members, ru.Nick)
	return nil
}

func (str *Storage) RoomHasUser(name string, c *cht.Client) (bool, int) {
	mr, exists := str.Rooms[name]
	if !exists {
		return false, 0
	}
	for _, mu := range mr.Members {
		if mu == UserId(c.Nick) {
			return true, 0
		}
	}
	return false, 0
}
