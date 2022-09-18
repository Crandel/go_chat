package sqlite

import (
	"math/rand"
	"time"

	cht "github.com/Crandel/go_chat/pkg/chatting"
	errs "github.com/Crandel/go_chat/pkg/errors"
)

type UserRoom struct {
	RoomName string `db:",rel=rooms"`
	UserID   string `db:",rel=users"`
}

func (str *Storage) WriteMessage(u cht.Client, r cht.Room, msg string) error {
	const op errs.Op = "sqlite.WriteMessage"
	room := Room{}
	err := str.db.Select(&room).Where("name = ?", r.Name).Do()
	if err != nil {
		return errs.NewError(
			op, errs.Info, "Room is not available", err)
	}
	rand.Seed(time.Now().UnixNano())
	id := rand.Int()
	message := Message{
		ID:       id,
		RoomName: room.Name,
		UserNick: u.GetNick(),
		Payload:  msg,
		Created:  time.Now(),
	}
	error := str.db.Insert(&message).Do()
	if error != nil {
		return errs.NewError(
			op, errs.Info, "", error)
	}

	return nil
}

func (str *Storage) ExcludeFromRoom(name string, u cht.Client) error {
	const op errs.Op = "sqlite.ExcludeFromRoom"
	// var mr Room
	// mr, exists := str.Rooms[roomName]
	// if !exists {
	// 	return errs.New(op, errs.Info, "No room with name "+roomName)
	// }
	// for i, ru := range mr.Members {
	// 	if string(ru) == *u.Nick {
	// 		mr.Members = append(mr.Members[:i], mr.Members[i+1:]...)
	// 	}
	// }
	// return nil
	// user, err := str.GetUser(reading.UserId(*u.Nick))
	// if err != nil {
	// 	return errs.New(op, errs.Info, "No user with nick: "+u.GetNick())
	// }
	// err = str.db.Delete()
	return nil
}

func (str *Storage) AddUserToRoom(name string, u cht.Client) error {
	return nil
}

func (str *Storage) RoomHasUser(name string, u cht.Client) bool {
	return false
}
