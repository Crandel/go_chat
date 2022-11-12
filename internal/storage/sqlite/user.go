package sqlite

import (
	"database/sql"
	"database/sql/driver"
	"time"

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
	Nick       string         `db:"nick,key"`
	Name       string         `db:"name"`
	SecondName sql.NullString `db:"second_name"`
	Email      sql.NullString `db:"email"`
	Password   string         `db:"password"`
	Token      string         `db:"token"`
	Role       Role           `db:"role"`
	Created    time.Time      `db:"created"`
}

func (*User) TableName() string {
	return USERS
}

func (u *User) GetSecondName() string {
	secondName := ""
	if u.SecondName.Valid {
		secondName = u.Email.String
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
		Name:       u.Name,
		SecondName: u.GetSecondName(),
	}
}
