package sqlite

import (
	"database/sql"

	errs "github.com/Crandel/go_chat/internal/errors"
	rdn "github.com/Crandel/go_chat/internal/reading"
)

func (str *Storage) ReadUsers() ([]rdn.User, error) {
	users := make([]User, 0)
	rdnUsers := make([]rdn.User, 0)
	err := str.db.Select(&users).Do()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		rdnUsers = append(rdnUsers, u.ConvertToReading())
	}
	return rdnUsers, nil
}

func (str *Storage) GetUser(ru rdn.UserId) (User, error) {
	const op errs.Op = "sqlite.GetUser"
	uid := string(ru)
	user := User{}
	err := str.db.Select(&user).Where("nick = ?", uid).Do()
	if err == sql.ErrNoRows {
		return User{}, errs.New(op, errs.Info, "No user with nick: "+uid)
	} else if err != nil {
		return User{}, errs.NewError(op, errs.Info, "Error with database connection", err)
	}
	return user, nil
}

func (str *Storage) ReadUser(ru rdn.UserId) (rdn.User, error) {
	const op errs.Op = "sqlite.ReadUser"
	user, err := str.GetUser(ru)
	if err != nil {
		return rdn.User{}, errs.New(op, errs.Info, "No user with nick: "+string(ru))
	}
	return user.ConvertToReading(), nil
}

func (str *Storage) ReadRooms() ([]rdn.Room, error) {
	rooms := make([]Room, 0)
	rdnRooms := make([]rdn.Room, 0)
	err := str.db.Select(&rooms).Do()
	if err != nil {
		return nil, err
	}
	for _, r := range rooms {
		rdnMessages := str.getRoomMessages(r.Name)
		rdnRooms = append(rdnRooms, r.ConvertToReading(rdnMessages))
	}
	return rdnRooms, nil
}

func (str *Storage) getRoom(id string) (Room, error) {
	const op errs.Op = "sqlite.getRoom"
	room := Room{}
	var newError error
	err := str.db.Select(&room).Where("name = ?", id).Do()
	if err == sql.ErrNoRows {
		newError = errs.New(op, errs.Info, "No room with name: "+id)
	} else if err != nil {
		newError = errs.NewError(op, errs.Info, "Error with database connection", err)
	}
	return room, newError

}

func (str *Storage) ReadRoom(id string) (rdn.Room, error) {
	const op errs.Op = "sqlite.ReadRoom"
	rdnRoom := rdn.Room{}
	room, err := str.getRoom(id)
	if err != nil {
		return rdnRoom, errs.NewError(op, errs.Info, "Can't get room", err)
	}

	rdnMessages := str.getRoomMessages(room.Name)
	return room.ConvertToReading(rdnMessages), nil
}

func (str *Storage) getRoomMessages(name string) map[rdn.UserId][]rdn.Message {
	messages := make([]Message, 0)
	rdnMessages := make(map[rdn.UserId][]rdn.Message, 0)
	users := str.getRoomUsers(name)
	for _, ur := range users {
		err := str.db.Select(messages).Where("user_room_id = ?", ur.id).Do()
		if err == nil {
			for _, m := range messages {
				rdnMessages[rdn.UserId(ur.userNick)] = append(rdnMessages[rdn.UserId(ur.userNick)], m.ConvertToReading(ur.userNick))
			}
		}

	}
	return rdnMessages
}
