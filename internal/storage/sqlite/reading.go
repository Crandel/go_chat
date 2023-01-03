package sqlite

import (
	"database/sql"

	lg "github.com/Crandel/go_chat/internal/logging"
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
	const op lg.Op = "sqlite.GetUser"
	uid := string(ru)
	user := User{}
	err := str.db.Select(&user).Where("nick = ?", uid).Do()
	if err == sql.ErrNoRows {
		return User{}, lg.New(op, "No user with nick: "+uid)
	} else if err != nil {
		return User{}, lg.NewError(op, "Error with database connection", err)
	}
	return user, nil
}

func (str *Storage) ReadUser(ru rdn.UserId) (rdn.User, error) {
	const op lg.Op = "sqlite.ReadUser"
	user, err := str.GetUser(ru)
	if err != nil {
		return rdn.User{}, lg.New(op, "No user with nick: "+string(ru))
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
	const op lg.Op = "sqlite.getRoom"
	room := Room{}
	var newError error
	err := str.db.Select(&room).Where("name = ?", id).Do()
	if err == sql.ErrNoRows {
		newError = lg.New(op, "No room with name: "+id)
	} else if err != nil {
		newError = lg.NewError(op, "Error with database connection", err)
	}
	return room, newError

}

func (str *Storage) ReadRoom(id string) (rdn.Room, error) {
	const op lg.Op = "sqlite.ReadRoom"
	rdnRoom := rdn.Room{}
	room, err := str.getRoom(id)
	if err != nil {
		return rdnRoom, lg.NewError(op, "Can't get room", err)
	}

	rdnMessages := str.getRoomMessages(room.Name)
	return room.ConvertToReading(rdnMessages), nil
}

func (str *Storage) getRoomMessages(name string) map[rdn.UserId][]rdn.Message {
	messages := make([]Message, 0)
	rdnMessages := make(map[rdn.UserId][]rdn.Message, 0)
	users := str.getRoomUsers(name)
	for _, ur := range users {
		err := str.db.Select(messages).Where("user_room_id = ?", ur.ID).Do()
		if err == nil {
			for _, m := range messages {
				rdnMessages[rdn.UserId(ur.UserNick)] = append(rdnMessages[rdn.UserId(ur.UserNick)], m.ConvertToReading(ur.UserNick))
			}
		}

	}
	return rdnMessages
}
