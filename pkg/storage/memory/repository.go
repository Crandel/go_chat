package memory

import (
	"errors"
	"fmt"
	"time"

	a "github.com/Crandel/go_chat/pkg/adding"
	err "github.com/Crandel/go_chat/pkg/errors"
	l "github.com/Crandel/go_chat/pkg/login"
	r "github.com/Crandel/go_chat/pkg/reading"
	s "github.com/Crandel/go_chat/pkg/signin"
)

type Storage struct {
	Users map[UserId]User
	Rooms map[string]Room
}

func NewStorage() Storage {
	return Storage{}
}

func (str *Storage) SigninUser(su s.User) (s.SigninResponse, error) {
	u := ConvertUserFromSigning(su)
	if str.Users == nil {
		str.Users = make(map[UserId]User)
	}
	str.Users[u.Email] = u
	return s.SigninResponse{Token: u.Token}, nil
}

func (str *Storage) LoginUser(lu l.User) (string, error) {
	const op err.Op = "memory.LoginUser"
	u, exists := str.Users[UserId(lu.Email)]
	if !exists {
		return "", err.New(op, err.Info, errors.New("No user with email: "+lu.Email))
	}
	if u.Password != lu.Password {
		return "", err.New(op, err.Info, errors.New("User with email"+lu.Email+"has wrong password"))
	}
	return u.Token, nil
}

func (str *Storage) AddRoom(ar a.Room) (string, error) {
	const op err.Op = "memory.AddRoom"
	if str.Rooms == nil {
		str.Rooms = make(map[string]Room)
	} else {
		_, exists := str.Rooms[ar.Name]
		if exists {
			return "", err.New(op, err.Info, errors.New(fmt.Sprintf("Room with name %s already exists", ar.Name)))
		}
	}
	messages := make(map[UserId][]Message)
	for _, user := range ar.Users {
		messages[UserId(user.ID)] = []Message{}
	}
	mr := Room{Name: ar.Name, Messages: messages, Created: time.Now()}
	str.Rooms[mr.Name] = mr
	return mr.Name, nil
}

func (str *Storage) ReadRooms() ([]r.Room, error) {
	const op err.Op = "memory.ReadRooms"
	var rooms = []r.Room{}
	for _, room := range str.Rooms {
		rooms = append(rooms, room.ConvertRoomToReading())
	}
	var not_found error
	if len(rooms) == 0 {
		not_found = err.New(op, err.Info, errors.New("No rooms are here"))
	}
	return rooms, not_found
}

func (str *Storage) ReadRoom(rid string) (r.Room, error) {
	const op err.Op = "memory.ReadRoom"
	room, exists := str.Rooms[rid]
	if !exists {
		return r.Room{}, err.New(op, err.Info, errors.New("No rooms with id "+rid))
	}
	return room.ConvertRoomToReading(), nil
}

func (str *Storage) ReadUsers() ([]r.User, error) {
	const op err.Op = "memory.ReadUsers"
	var users = []r.User{}
	for _, u := range str.Users {
		users = append(users, u.ConvertUserToReading())
	}
	var not_found error
	if len(users) == 0 {
		not_found = err.New(op, err.Info, errors.New("No users are here"))
	}
	return users, not_found
}

func (str *Storage) ReadUser(uid r.UserId) (r.User, error) {
	const op err.Op = "memory.ReadUser"
	umid := ConvertUserIdFromReading(uid)
	s_user, exists := str.Users[umid]
	if !exists {
		return r.User{}, err.New(op, err.Info, errors.New("No user with id"+string(uid)))
	}
	return s_user.ConvertUserToReading(), nil
}
