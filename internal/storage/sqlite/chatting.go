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
	id       int    `db:"id,key,auto"`
	roomName string `db:",rel=rooms"`
	userNick string `db:",rel=users"`
}

func (*UserRoom) TableName() string {
	return USER_ROOMS
}

func (str *Storage) WriteMessage(u *cht.Client, r *cht.Room, msg string) error {
	const op errs.Op = "sqlite.WriteMessage"
	exists, id := str.RoomHasUser(r.Name, u)
	if !exists {
		return errs.New(
			op, errs.Info, "User "+*u.Nick+"is not in room "+r.Name)
	}
	rand.Seed(time.Now().UnixNano())
	message := Message{
		UserRoomID: id,
		Payload:    msg,
		Created:    time.Now(),
	}
	error := str.db.Insert(&message).Do()
	if error != nil {
		return errs.NewError(
			op, errs.Info, "", error)
	}

	return nil
}

func (str *Storage) ExcludeFromRoom(name string, u *cht.Client) error {
	const op errs.Op = "sqlite.ExcludeFromRoom"
	exists, id := str.RoomHasUser(name, u)

	if !exists {
		return errs.New(
			op, errs.Info, "User "+*u.Nick+"is not in room "+name)
	}
	_, err := str.db.DeleteFrom(
		USER_ROOMS).Where("id = ?", id).Do()
	if err != nil {
		return errs.New(
			op, errs.Info, "User "+*u.Nick+"is not in room "+name)
	}
	return nil
}

func (str *Storage) AddUserToRoom(name string, c *cht.Client) error {
	const op errs.Op = "sqlite.AddUserToRoom"
	exists, _ := str.RoomHasUser(name, c)
	if exists {
		return errs.New(
			op, errs.Info, "User "+*c.Nick+" is already in a room "+name)
	}

	userInRoom := UserRoom{
		roomName: name,
		userNick: *c.Nick,
	}
	error := str.db.Insert(&userInRoom).Do()
	if error != nil {
		return errs.NewError(
			op, errs.Info, "", error)
	}
	return nil
}

func (str *Storage) getRoomUsers(name string) []UserRoom {
	var usersInRoom []UserRoom
	_ = str.db.Select(&usersInRoom).WhereQ(
		godb.And(
			godb.Q("room_name = ?", name),
		)).Do()
	return usersInRoom
}

func (str *Storage) RoomHasUser(name string, c *cht.Client) (bool, int) {
	var userInRoom UserRoom
	_ = str.db.Select(&userInRoom).WhereQ(
		godb.And(
			godb.Q("room_name = ?", name),
			godb.Q("user_nick = ?", c.Nick),
		)).Do()
	if userInRoom != (UserRoom{}) {
		return true, userInRoom.id
	}
	return false, 0
}
