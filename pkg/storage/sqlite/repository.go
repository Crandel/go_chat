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
			op, errs.Info, fmt.Sprintf("Insert user with email %s is not possible", u.Email), error)
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
	} else {
		return user.Token, errs.NewError(op, errs.Info, "Internal error", error)
	}
}

func (str *Storage) AddRoom(ar a.Room) (string, error) {
	return "", nil
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
