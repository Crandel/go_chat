package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	a "github.com/Crandel/go_chat/pkg/adding"
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
	err := str.db.Insert(&su).Do()
	if err != nil {
		return s.SigninResponse{}, fmt.Errorf("Insert user with email %s is not possible, due to error: %w", u.Email, err)
	}
	return s.SigninResponse{Token: token}, nil
}

func (str *Storage) LoginUser(lu l.User) (string, error) {
	user := User{}
	err := str.db.Select(&user).WhereQ(
		godb.And(
			godb.Q("email = ?", lu.Email),
			godb.Q("password = ? ", lu.Password),
		),
	).Do()
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("No user with email: " + lu.Email)
	} else {
		return user.Token, err
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
