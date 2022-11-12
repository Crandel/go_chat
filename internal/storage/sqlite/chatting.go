package sqlite

import (
	"math/rand"
	"time"

	cht "github.com/Crandel/go_chat/internal/chatting"
	errs "github.com/Crandel/go_chat/internal/errors"
	"github.com/samonzeweb/godb"
)

const USER_ROOMS = "user_rooms"

type UserRoom struct {
	roomName string `db:",rel=rooms"`
	userNick string `db:",rel=users"`
}

func (*UserRoom) TableName() string {
	return USER_ROOMS
}

func (str *Storage) WriteMessage(u cht.Client, r cht.Room, msg string) error {
	const op errs.Op = "sqlite.WriteMessage"
	if str.RoomHasUser(r.Name, u) {
		return errs.New(
			op, errs.Info, "User "+*u.Nick+"is not in room "+r.Name)
	}
	rand.Seed(time.Now().UnixNano())
	id := rand.Int()
	message := Message{
		ID:       id,
		RoomName: r.Name,
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
	if !str.RoomHasUser(name, u) {
		return errs.New(
			op, errs.Info, "User "+*u.Nick+"is not in room "+name)
	}
	_, err := str.db.DeleteFrom(
		USER_ROOMS).WhereQ(
		godb.And(
			godb.Q("room_name = ?", name),
			godb.Q("", u.Nick))).Do()
	if err != nil {
		return errs.New(
			op, errs.Info, "User "+*u.Nick+"is not in room "+name)
	}
	return nil
}

func (str *Storage) AddUserToRoom(name string, u cht.Client) error {
	const op errs.Op = "sqlite.AddUserToRoom"
	if str.RoomHasUser(name, u) {
		return errs.New(
			op, errs.Info, "User "+*u.Nick+" is already in a room "+name)
	}

	userInRoom := UserRoom{
		roomName: name,
		userNick: *u.Nick,
	}
	error := str.db.Insert(&userInRoom).Do()
	if error != nil {
		return errs.NewError(
			op, errs.Info, "", error)
	}
	return nil
}

func (str *Storage) RoomHasUser(name string, u cht.Client) bool {
	var userInRoom UserRoom
	_ = str.db.Select(&userInRoom).WhereQ(
		godb.And(
			godb.Q("room_name = ?", name),
			godb.Q("user_nick = ?", u.Nick),
		)).Do()
	if userInRoom != (UserRoom{}) {
		return true
	}
	return false
}
