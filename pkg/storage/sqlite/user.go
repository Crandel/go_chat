package sqlite

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/Crandel/go_chat/pkg/login"
	"github.com/Crandel/go_chat/pkg/signin"
	"github.com/google/uuid"
	"github.com/samonzeweb/godb"
)

type Role string

const (
	Member Role = "Member"
	Admin  Role = "Admin"
)

func (r *Role) Scan(value interface{}) error { *r = Role(value.(string)); return nil }
func (r Role) Value() (driver.Value, error)  { return driver.Value(string(r)), nil }

type User struct {
	ID         string    `db:"id"`
	Name       string    `db:"name"`
	SecondName string    `db:"second_name"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	Token      string    `db:"token"`
	Role       Role      `db:"role"`
	Created    time.Time `db:"created"`
}

func (*User) TableName() string {
	return "users"
}

func (s *Storage) SigninUser(u signin.User) (signin.SigninResponse, error) {
	id := uuid.New().String()
	token := uuid.New().String()
	su := User{
		ID:         id,
		Name:       u.Name,
		SecondName: u.SecondName,
		Email:      u.Email,
		Password:   u.Password,
		Token:      token,
		Role:       Member,
		Created:    time.Now(),
	}
	s.db.InsertInto("users").
		Columns("name", "second_name", "email", "password", "token", "role", "created").
		Values(su.Name, su.SecondName, su.Email, su.Password, su.Token, su.Role, su.Created).
		Returning("id").
		DoWithReturning(&su)
	return signin.SigninResponse{Id: fmt.Sprint(su.ID), Token: token}, nil
}

func (s *Storage) LoginUser(lu login.User) (string, error) {
	user := User{}
	err := s.db.Select(&user).WhereQ(
		godb.And(
			godb.Q("email = ?", lu.Email),
			godb.Q("password = ? ", lu.Password),
		),
	).Do()
	if err == sql.ErrNoRows {
		return "", errors.New("No user with email: " + lu.Email)
	} else {
		return user.Token, err
	}
}
