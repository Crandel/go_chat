package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	a "github.com/Crandel/go_chat/pkg/adding"
	errs "github.com/Crandel/go_chat/pkg/errors"
	l "github.com/Crandel/go_chat/pkg/login"
	r "github.com/Crandel/go_chat/pkg/reading"
	s "github.com/Crandel/go_chat/pkg/signin"
	"github.com/google/uuid"
	"github.com/samonzeweb/godb"
)

type Storage struct {
	db *godb.DB
}

func NewStorage(db *godb.DB) Storage {
	return Storage{db}
}

func (str *Storage) SigninUser(u s.User) (s.SigninResponse, error) {
	const op errs.Op = "sqlite.SigninUser"
	token := uuid.New().String()
	su := User{
		Name:       u.Name,
		SecondName: u.SecondName,
		Email:      u.Email,
		Password:   u.Password,
		Token:      token,
		Role:       Member,
		Created:    time.Now(),
	}
	error := str.db.Insert(&su).Do()
	if error != nil {
		return s.SigninResponse{}, errs.NewError(
			op, errs.Info, fmt.Sprintf("User with email %s already exists", u.Email), error)
	}
	return s.SigninResponse{Token: token}, nil
}

func (str *Storage) LoginUser(lu l.User) (string, error) {
	const op errs.Op = "sqlite.LoginUser"
	user := User{}
	error := str.db.Select(&user).WhereQ(
		godb.And(
			godb.Q("email = ?", lu.Email),
			godb.Q("password = ? ", lu.Password),
		),
	).Do()
	if error == sql.ErrNoRows {
		return "", errs.New(op, errs.Info, "No user with email: "+lu.Email)
	} else if error != nil {
		return "", errs.NewError(op, errs.Info, "Internal error", error)
	}
	return user.Token, nil
}

func (str *Storage) AddRoom(ar a.Room) (string, []error) {
	const op errs.Op = "sqlite.AddRoom"
	users := []UserMessage{}
	list_errors := []error{}
	for _, au := range ar.Users {
		su := User{}
		error := str.db.Select(&su).Where("id = ?", au.ID).Do()
		if error != nil {
			list_errors = append(
				list_errors,
				errs.NewError(op, errs.Info, "User with ID"+au.ID+"not found", error))
		} else {
			users = append(users, UserMessage{User: su, Messages: []Message{}})
		}
	}
	res_str := ""
	if len(list_errors) == 0 {
		id := uuid.New().String()
		room := Room{id, users}
		error := str.db.Insert(&room).Do()
		if error != nil {
			res_str = room.ID
		} else {
			list_errors = append(
				list_errors,
				errs.NewError(op, errs.Info, "Failed to create room", error))
		}
	}
	return res_str, list_errors
}

func (str *Storage) ReadUsers() ([]r.User, error) {
	return []r.User{}, nil
}

func (str *Storage) ReadUser(r.UserId) (r.User, error) {
	return r.User{}, nil
}

func (str *Storage) ReadRooms() ([]r.Room, error) {
	return []r.Room{}, nil
}

func (str *Storage) ReadRoom(string) (r.Room, error) {
	return r.Room{}, nil
}
