package memory

import (
	"fmt"
	"sync"
	"time"

	a "github.com/Crandel/go_chat/pkg/adding"
	errs "github.com/Crandel/go_chat/pkg/errors"
	l "github.com/Crandel/go_chat/pkg/login"
	r "github.com/Crandel/go_chat/pkg/reading"
	s "github.com/Crandel/go_chat/pkg/signin"
)

type Storage struct {
	Users map[UserId]User
	Rooms map[string]Room
	*sync.RWMutex
}

func NewStorage() Storage {
	return Storage{
		Users:   map[UserId]User{},
		Rooms:   map[string]Room{},
		RWMutex: &sync.RWMutex{},
	}
}

func (str *Storage) SigninUser(su s.User) (s.SigninResponse, error) {
	u := ConvertUserFromSigning(su)
	if str.Users == nil {
		str.Users = make(map[UserId]User)
	}
	str.Lock()
	str.Users[u.Email] = u
	str.Unlock()
	return s.SigninResponse{Token: u.Token}, nil
}

func (str *Storage) LoginUser(lu l.User) (string, error) {
	const op errs.Op = "memory.LoginUser"
	str.RLock()
	u, exists := str.Users[UserId(lu.Email)]
	str.RUnlock()
	if !exists {
		return "", errs.New(op, errs.Info, "No user with email: "+lu.Email)
	}
	if u.Password != lu.Password {
		return "", errs.New(op, errs.Info, "User with email"+lu.Email+"has wrong password")
	}
	return u.Token, nil
}

func (str *Storage) AddRoom(ar a.Room) (string, []error) {
	const op errs.Op = "memory.AddRoom"
	var error_list []error
	str.Lock()
	if str.Rooms == nil {
		str.Rooms = make(map[string]Room)
	} else {
		_, exists := str.Rooms[ar.Name]
		if exists {
			error_list = append(
				error_list,
				errs.New(op, errs.Info, fmt.Sprintf("Room with name %s already exists", ar.Name)))
			return "", error_list
		}
	}
	str.Unlock()
	messages := make(map[UserId][]Message)
	for _, user := range ar.Users {
		uid := UserId(user.ID)
		str.RLock()
		u, exists := str.Users[uid]
		str.RUnlock()
		if !exists {
			error_list = append(error_list, errs.New(op, errs.Info, "User with id "+user.ID+"does not exists"))
		} else {
			messages[u.Email] = []Message{}
		}
	}
	mr := Room{Name: ar.Name, Messages: messages, Created: time.Now()}
	str.Lock()
	str.Rooms[mr.Name] = mr
	str.Unlock()
	return mr.Name, error_list
}

func (str *Storage) ReadRooms() ([]r.Room, error) {
	const op errs.Op = "memory.ReadRooms"
	var rooms = []r.Room{}
	str.RLock()
	for _, room := range str.Rooms {
		rooms = append(rooms, room.ConvertRoomToReading())
	}
	str.RUnlock()
	var not_found error
	if len(rooms) == 0 {
		not_found = errs.New(op, errs.Info, "No rooms are here")
	}
	return rooms, not_found
}

func (str *Storage) ReadRoom(rid string) (r.Room, error) {
	const op errs.Op = "memory.ReadRoom"
	str.RLock()
	room, exists := str.Rooms[rid]
	str.RUnlock()
	if !exists {
		return r.Room{}, errs.New(op, errs.Info, "No rooms with id "+rid)
	}
	return room.ConvertRoomToReading(), nil
}

func (str *Storage) ReadUsers() ([]r.User, error) {
	const op errs.Op = "memory.ReadUsers"
	var users = []r.User{}
	str.RLock()
	for _, u := range str.Users {
		users = append(users, u.ConvertUserToReading())
	}
	str.RUnlock()
	var not_found error
	if len(users) == 0 {
		not_found = errs.New(op, errs.Info, "No users are here")
	}
	return users, not_found
}

func (str *Storage) ReadUser(uid r.UserId) (r.User, error) {
	const op errs.Op = "memory.ReadUser"
	umid := ConvertUserIdFromReading(uid)
	str.RLock()
	s_user, exists := str.Users[umid]
	str.RUnlock()
	if !exists {
		return r.User{}, errs.New(op, errs.Info, "No user with id "+string(uid))
	}
	return s_user.ConvertUserToReading(), nil
}
