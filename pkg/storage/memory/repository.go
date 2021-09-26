package memory

import (
	"fmt"
	"sync"
	"time"

	add "github.com/Crandel/go_chat/pkg/adding"
	"github.com/Crandel/go_chat/pkg/auth"
	errs "github.com/Crandel/go_chat/pkg/errors"
	rdn "github.com/Crandel/go_chat/pkg/reading"
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

func FilledStorage(users map[UserId]User, rooms map[string]Room) Storage {
	return Storage{
		Users:   users,
		Rooms:   rooms,
		RWMutex: &sync.RWMutex{},
	}
}

func (str *Storage) SigninUser(su auth.SigninUser) (string, error) {
	const op errs.Op = "memory.Signin"
	u := ConvertUserFromSigning(su)
	if str.Users == nil {
		str.Users = make(map[UserId]User)
	}
	str.Lock()
	_, exists := str.Users[u.Email]
	if exists {
		return "", errs.New(
			op, errs.Info, fmt.Sprintf("User with email: '%s' exists", u.Email))
	}
	str.Users[u.Email] = u
	str.Unlock()
	return u.Token, nil
}

func (str *Storage) LoginUser(lu auth.LoginUser) (string, error) {
	const op errs.Op = "memory.LoginUser"
	str.RLock()
	u, exists := str.Users[UserId(lu.Email)]
	str.RUnlock()
	if !exists {
		return "", errs.New(
			op, errs.Info, fmt.Sprintf("No user with email '%s'", lu.Email))
	}
	if u.Password != lu.Password {
		return "", errs.New(
			op, errs.Info, fmt.Sprintf("User with email '%s' has wrong password", lu.Email))
	}
	return u.Token, nil
}

func (str *Storage) AddRoom(ar add.Room) (string, []error) {
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
				errs.New(
					op, errs.Info, fmt.Sprintf("Room with name %s already exists", ar.Name)))
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
			error_list = append(error_list, errs.New(
				op, errs.Info, fmt.Sprintf("User with email '%s' does not exists", user.ID)))
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

func (str *Storage) ReadRooms() ([]rdn.Room, error) {
	const op errs.Op = "memory.ReadRooms"
	var rooms = []rdn.Room{}
	str.RLock()
	for _, room := range str.Rooms {
		rooms = append(rooms, room.ConvertRoomToReading())
	}
	str.RUnlock()
	return rooms, nil
}

func (str *Storage) ReadRoom(rid string) (rdn.Room, error) {
	const op errs.Op = "memory.ReadRoom"
	str.RLock()
	room, exists := str.Rooms[rid]
	str.RUnlock()
	if !exists {
		return rdn.Room{}, errs.New(op, errs.Info, "No rooms with id "+rid)
	}
	return room.ConvertRoomToReading(), nil
}

func (str *Storage) ReadUsers() ([]rdn.User, error) {
	const op errs.Op = "memory.ReadUsers"
	var users = []rdn.User{}
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

func (str *Storage) ReadUser(uid rdn.UserId) (rdn.User, error) {
	const op errs.Op = "memory.ReadUser"
	umid := ConvertUserIdFromReading(uid)
	str.RLock()
	s_user, exists := str.Users[umid]
	str.RUnlock()
	if !exists {
		return rdn.User{}, errs.New(op, errs.Info, "No user with id "+string(uid))
	}
	return s_user.ConvertUserToReading(), nil
}
