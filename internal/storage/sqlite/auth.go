package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Crandel/go_chat/internal/auth"
	errs "github.com/Crandel/go_chat/internal/errors"
	"github.com/google/uuid"
	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/types"
)

func (str *Storage) SigninUser(u auth.SigninUser) (string, error) {
	const op errs.Op = "sqlite.SigninUser"
	token := uuid.New().String()
	su := User{
		Name:       u.Name,
		SecondName: types.NullStringFrom(*u.SecondName),
		Email:      types.NullStringFrom(*u.Email),
		Password:   u.Password,
		Token:      token,
		Role:       Member,
		Created:    time.Now(),
	}
	err := str.db.Insert(&su).Do()
	if err != nil {
		return "", errs.NewError(
			op, errs.Info, fmt.Sprintf("User with nick %s already exists", u.Nick), err)
	}
	return token, nil
}

func (str *Storage) LoginUser(lu auth.LoginUser) (string, error) {
	const op errs.Op = "sqlite.LoginUser"
	user := User{}
	err := str.db.Select(&user).WhereQ(
		godb.And(
			godb.Q("nick = ?", lu.Nick),
			godb.Q("password = ? ", lu.Password),
		),
	).Do()
	if err == sql.ErrNoRows {
		return "", errs.New(op, errs.Info, "No user with nick: "+lu.Nick)
	} else if err != nil {
		return "", errs.NewError(op, errs.Info, "Internal error", err)
	}

	return user.Token, nil
}
