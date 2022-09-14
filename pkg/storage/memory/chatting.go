package memory

import (
	"math/rand"
	"time"

	cht "github.com/Crandel/go_chat/pkg/chatting"
	errs "github.com/Crandel/go_chat/pkg/errors"
)

func (str *Storage) WriteMessage(u cht.Client, r cht.Room, msg string) error {
	const op errs.Op = "memory.WriteMessage"
	rand.Seed(time.Now().UnixNano())
	id := rand.Int()
	str.Messages[id] = Message{
		ID:       id,
		UserId:   UserId(*u.Nick),
		RoomName: r.Name,
		Payload:  msg,
		Created:  time.Now(),
	}
	return nil
}

func (str *Storage) ExcludeFromRoom(roomName string, u cht.Client) error {
	const op errs.Op = "memory.ExcludeFromRoom"
	var mr Room
	mr, exists := str.Rooms[roomName]
	if !exists {
		return errs.New(op, errs.Info, "No room with name "+roomName)
	}
	for i, ru := range mr.Members {
		if string(ru) == *u.Nick {
			mr.Members = append(mr.Members[:i], mr.Members[i+1:]...)
		}
	}
	return nil
}

func (str *Storage) AddUserToRoom(roomName string, u cht.Client) error {
	const op errs.Op = "memory.AddUserToRoom"
	var mr Room
	mr, exists := str.Rooms[roomName]
	if !exists {
		mrname, error := str.AddRoom(roomName)
		if error != nil {
			return error
		}
		mr = str.Rooms[mrname]
	}
	ru, exists := str.Users[UserId(*u.Nick)]
	if !exists {
		return errs.New(op, errs.Info, "No user with id "+*u.Nick)
	}
	mr.Members = append(mr.Members, ru.Email)
	return nil
}

func (str *Storage) RoomHasUser(roomName string, cu cht.Client) bool {
	const op errs.Op = "memory.RoomHasUser"
	mr, exists := str.Rooms[roomName]
	if !exists {
		return false
	}
	for _, mu := range mr.Members {
		if mu == UserId(*cu.Nick) {
			return true
		}
	}
	return false
}
