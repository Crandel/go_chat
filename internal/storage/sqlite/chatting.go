package sqlite

import (
	"math/rand"
	"time"

	cht "github.com/Crandel/go_chat/internal/chatting"
	lg "github.com/Crandel/go_chat/internal/logging"
	"github.com/Crandel/go_chat/internal/reading"
	"github.com/samonzeweb/godb"
)

const USER_ROOMS = "user_rooms"

type UserRoom struct {
	RoomName string `db:"room_name"`
	UserNick string `db:"user_nick"`
	ID       int    `db:"id,key,auto"`
}

func (*UserRoom) TableName() string {
	return USER_ROOMS
}

func (str *Storage) WriteMessage(u *cht.Client, r *cht.Room, msg string) error {
	const op lg.Op = "sqlite.WriteMessage"
	exists, id := str.RoomHasUser(r.Name, u)
	if !exists {
		return lg.New(
			op, "User "+u.Nick+"is not in room "+r.Name)
	}
	rand.Seed(time.Now().UnixNano())
	message := Message{
		UserRoomID: id,
		Payload:    msg,
		Created:    time.Now(),
	}
	error := str.db.Insert(&message).Do()
	if error != nil {
		return lg.NewError(
			op, "", error)
	}

	return nil
}

func (str *Storage) ExcludeFromRoom(name string, u *cht.Client) error {
	const op lg.Op = "sqlite.ExcludeFromRoom"
	exists, id := str.RoomHasUser(name, u)

	if !exists {
		return lg.New(
			op, "User "+u.Nick+"is not in room "+name)
	}
	_, err := str.db.DeleteFrom(
		USER_ROOMS).Where("id = ?", id).Do()
	if err != nil {
		return lg.New(
			op, "User "+u.Nick+"is not in room "+name)
	}
	return nil
}

func (str *Storage) AddUserToRoom(name string, c *cht.Client) error {
	const op lg.Op = "sqlite.AddUserToRoom"
	exists, _ := str.RoomHasUser(name, c)
	if exists {
		return nil
	}

	user, error := str.GetUser(reading.UserId(c.Nick))
	if error != nil {
		return lg.NewError(
			op, "User not found", error)
	}
	room, error := str.getRoom(name)
	if error != nil {
		return lg.NewError(
			op, "Room not found", error)
	}

	userInRoom := UserRoom{
		RoomName: room.Name,
		UserNick: user.Nick,
	}

	log.Debugln(userInRoom)
	error = str.db.Insert(&userInRoom).Do()
	if error != nil {
		return lg.NewError(
			op, "", error)
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
		return true, userInRoom.ID
	}
	return false, 0
}
