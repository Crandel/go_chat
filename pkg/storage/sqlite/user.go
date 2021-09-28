package sqlite

import (
	"database/sql/driver"
	"time"

	rdn "github.com/Crandel/go_chat/pkg/reading"
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
	ID         string    `db:"id,key"`
	Name       string    `db:"name"`
	SecondName string    `db:"second_name"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	Token      string    `db:"token"`
	Role       Role      `db:"role"`
	Created    time.Time `db:"created"`
}

func (*User) TableName() string {
	return USERS
}

func (u *User) ConvertToReading() rdn.User {
	return rdn.User{
		ID:         rdn.UserId(u.ID),
		Email:      u.Email,
		Name:       u.Name,
		SecondName: u.SecondName,
	}
}
