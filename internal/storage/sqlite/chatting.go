package sqlite

import (
	"math/rand"
	"time"

	cht "github.com/Crandel/go_chat/internal/chatting"
	errs "github.com/Crandel/go_chat/internal/errors"
	"github.com/Crandel/go_chat/internal/reading"
	"github.com/samonzeweb/godb"
)

const USER_ROOMS = "user_rooms"

type UserRoom struct {
	roomName Room `db:",rel=rooms"`
	userNick User `db:",rel=users"`
	id       int  `db:"id,key,auto"`
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

	user, error := str.GetUser(reading.UserId(c.GetNick()))
	if error != nil {
		return errs.NewError(
			op, errs.Info, "User not found", error)
	}
	room, error := str.getRoom(name)
	if error != nil {
		return errs.NewError(
			op, errs.Info, "Room not found", error)
	}

	userInRoom := UserRoom{
		roomName: room,
		userNick: user,
	}

	error = str.db.Insert(&userInRoom).Do()
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
