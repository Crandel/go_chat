package sqlite

import (
	"database/sql"

	errs "github.com/Crandel/go_chat/pkg/errors"
	rdn "github.com/Crandel/go_chat/pkg/reading"
)

func (str *Storage) ReadUsers() ([]rdn.User, error) {
	const op errs.Op = "sqlite.LoginUser"
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

func (str *Storage) ReadUser(ru rdn.UserId) (rdn.User, error) {
	const op errs.Op = "sqlite.ReadUser"
	uid := string(ru)
	user := User{}
	err := str.db.Select(&user).Where("email = ?", uid).Do()
	if err == sql.ErrNoRows {
		return rdn.User{}, errs.New(op, errs.Info, "No user with id: "+uid)
	} else if err != nil {
		return rdn.User{}, errs.NewError(op, errs.Info, "Error with database connection", err)
	}
	return user.ConvertToReading(), nil
}

func (str *Storage) ReadRooms() ([]rdn.Room, error) {
	const op errs.Op = "sqlite.ReadRooms"
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

func (str *Storage) ReadRoom(id string) (rdn.Room, error) {
	const op errs.Op = "sqlite.ReadRoom"
	room := Room{}
	rdnRoom := rdn.Room{}
	err := str.db.Select(&room).Where("name = ?", id).Do()
	if err == sql.ErrNoRows {
		return rdnRoom, errs.New(op, errs.Info, "No room with name: "+id)
	} else if err != nil {
		return rdnRoom, errs.NewError(op, errs.Info, "Error with database connection", err)
	}
	rdnMessages := str.getRoomMessages(room.Name)
	return room.ConvertToReading(rdnMessages), nil
}

func (str *Storage) getRoomMessages(name string) map[rdn.UserId][]rdn.Message {
	messages := make([]Message, 0)
	rdnMessages := make(map[rdn.UserId][]rdn.Message, 0)
	err := str.db.Select(messages).Where("room_name = ?", name).Do()
	if err == nil {
		for _, m := range messages {
			rdnMessages[rdn.UserId(m.UserID)] = append(rdnMessages[rdn.UserId(m.UserID)], m.ConvertToReading())
		}
	}
	return rdnMessages
}
