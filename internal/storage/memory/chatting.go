package memory

import (
	"math/rand"
	"time"

	cht "github.com/Crandel/go_chat/internal/chatting"
	errs "github.com/Crandel/go_chat/internal/errors"
)

func (str *Storage) WriteMessage(u *cht.Client, r *cht.Room, msg string) error {
	const op errs.Op = "memory.WriteMessage"
	rand.Seed(time.Now().UnixNano())
	id := rand.Int()
	str.Messages[id] = Message{
		ID:       id,
		UserID:   UserId(u.GetNick()),
		RoomName: r.Name,
		Payload:  msg,
		Created:  time.Now(),
	}
	return nil
}

func (str *Storage) ExcludeFromRoom(roomName string, u *cht.Client) error {
	const op errs.Op = "memory.ExcludeFromRoom"
	var mr Room
	mr, exists := str.Rooms[roomName]
	if !exists {
		return errs.New(op, errs.Info, "No room with name "+roomName)
	}
	for i, ru := range mr.Members {
		if string(ru) == u.GetNick() {
			mr.Members = append(mr.Members[:i], mr.Members[i+1:]...)
		}
	}
	return nil
}

func (str *Storage) AddUserToRoom(name string, c *cht.Client) error {
	const op errs.Op = "memory.AddUserToRoom"
	var mr Room
	mr, exists := str.Rooms[name]
	if !exists {
		mrname, error := str.AddRoom(name)
		if error != nil {
			return error
		}
		mr = str.Rooms[mrname]
	}
	ru, exists := str.Users[UserId(c.GetNick())]
	if !exists {
		return errs.New(op, errs.Info, "No user with id "+c.GetNick())
	}
	mr.Members = append(mr.Members, ru.Nick)
	return nil
}

func (str *Storage) RoomHasUser(name string, c *cht.Client) (bool, int) {
	const op errs.Op = "memory.RoomHasUser"
	mr, exists := str.Rooms[name]
	if !exists {
		return false, 0
	}
	for _, mu := range mr.Members {
		if mu == UserId(c.GetNick()) {
			return true, 0
		}
	}
	return false, 0
}