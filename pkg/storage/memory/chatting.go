package memory

import (
	cht "github.com/Crandel/go_chat/pkg/chatting"
	errs "github.com/Crandel/go_chat/pkg/errors"
)

func (str *Storage) WriteMessage(u cht.User, r cht.Room, msg string) error {
	return nil
}

func (str *Storage) ExcludeFromRoom(roomName string, u cht.User) error {
	return nil
}

func (str *Storage) AddUserToRoom(roomName string, u cht.User) error {
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
	ru, exists := str.Users[UserId(u.Nick)]
	if !exists {
		return errs.New(op, errs.Info, "No user with id "+u.Nick)
	}
	mr.Members = append(mr.Members, ru.Email)
	return nil
}

func (str *Storage) RoomHasUser(roomName string, cu cht.User) bool {
	const op errs.Op = "memory.RoomHasUser"
	mr, exists := str.Rooms[roomName]
	if !exists {
		return false
	}
	for _, mu := range mr.Members {
		if mu == UserId(cu.Nick) {
			return true
		}
	}
	return false
}
