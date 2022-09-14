package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Crandel/go_chat/pkg/auth"
	errs "github.com/Crandel/go_chat/pkg/errors"
	"github.com/google/uuid"
	"github.com/samonzeweb/godb"
)

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
