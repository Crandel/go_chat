package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	add "github.com/Crandel/go_chat/pkg/adding"
	"github.com/Crandel/go_chat/pkg/auth"
	errs "github.com/Crandel/go_chat/pkg/errors"
	rdn "github.com/Crandel/go_chat/pkg/reading"
	"github.com/google/uuid"
	"github.com/samonzeweb/godb"
)

type Storage struct {
	db *godb.DB
}

func NewStorage(db *godb.DB) Storage {
	return Storage{db}
}

func (str *Storage) SigninUser(u auth.SigninUser) (string, error) {
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
		return "", errs.NewError(
			op, errs.Info, fmt.Sprintf("User with email %s already exists", u.Email), error)
	}
	return token, nil
}

func (str *Storage) LoginUser(lu auth.LoginUser) (string, error) {
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

func (str *Storage) AddRoom(ar add.Room) (string, []error) {
	const op errs.Op = "sqlite.AddRoom"
	room := Room{}
	err := str.db.Select(&room).Where("name = ?", ar.Name).Do()
	if err == nil && &room == nil {
		return "", []error{errs.New(op, errs.Info, "Room already exists")}
	}
	list_errors := []error{}
	id := uuid.New().String()
	room.Name = id
	room.Created = time.Now()
	error := str.db.Insert(&room).Do()
	res_str := ""
	if error != nil {
		res_str = room.Name
	} else {
		list_errors = append(
			list_errors,
			errs.NewError(op, errs.Info, "Failed to create room", error))
	}
	return res_str, list_errors
}

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
	return []rdn.Room{}, nil
}

func (str *Storage) ReadRoom(id string) (rdn.Room, error) {
	return rdn.Room{}, nil
}
