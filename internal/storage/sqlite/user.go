package sqlite

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/Crandel/go_chat/internal/auth"
	rdn "github.com/Crandel/go_chat/internal/reading"
)

const USERS = "users"

type Role string

const (
	Member Role = "Member"
	Admin  Role = "Admin"
)

func (r *Role) Scan(value interface{}) error {
	*r = Role(value.(string))
	return nil
}
func (r Role) Value() (driver.Value, error) {
	return driver.Value(string(r)), nil
}

type User struct {
	Created    time.Time      `db:"created"`
	Nick       string         `db:"nick,key"`
	Token      string         `db:"token"`
	Role       Role           `db:"role"`
	Name       sql.NullString `db:"name"`
	SecondName sql.NullString `db:"second_name"`
	Email      sql.NullString `db:"email"`
}

func (*User) TableName() string {
	return USERS
}

func (u *User) GetName() string {
	name := ""
	if u.Name.Valid {
		name = u.Name.String
	}
	return name
}
func (u *User) GetSecondName() string {
	secondName := ""
	if u.SecondName.Valid {
		secondName = u.SecondName.String
	}
	return secondName
}

func (u *User) GetEmail() string {
	email := ""
	if u.Email.Valid {
		email = u.Email.String
	}
	return email
}

func (u *User) ConvertToReading() rdn.User {
	return rdn.User{
		Nick:       rdn.UserId(u.Nick),
		Email:      u.GetEmail(),
		Name:       u.GetName(),
		SecondName: u.GetSecondName(),
	}
}

func (u *User) ConvertUserToAuth() auth.AuthUser {
	return auth.AuthUser{
		Nick:  string(u.Nick),
		Token: u.Token,
	}
}
