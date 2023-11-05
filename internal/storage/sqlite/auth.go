package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Crandel/go_chat/internal/auth"
	lg "github.com/Crandel/go_chat/internal/logging"
	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/types"
)

func (str *Storage) SigninUser(u auth.SigninUser) (string, error) {
	const op lg.Stk = "sqlite.SigninUser"

	token := auth.MakeToken(u.Nick, u.Password)
	su := User{
		Nick:  u.Nick,
		Token: token,
		Role:  Member,
	}
	if u.Name != nil {
		su.Name = types.NullStringFrom(*u.Name)
	}
	if u.SecondName != nil {
		su.SecondName = types.NullStringFrom(*u.SecondName)
	}
	if u.Email != nil {
		su.Email = types.NullStringFrom(*u.Email)
	}

	err := str.db.Insert(&su).Do()
	if err != nil {
		return "", lg.NewError(
			op, fmt.Sprintf("User with nick %s already exists", u.Nick), err)
	}
	return token, nil
}

func (str *Storage) LoginUser(lu auth.LoginUser) (string, error) {
	const op lg.Stk = "sqlite.LoginUser"

	user := User{}
	token := auth.MakeToken(lu.Nick, lu.Password)
	slog.Debug("Token %s; Nick %s ", token, lu.Nick)
	query := str.db.Select(&user).WhereQ(
		godb.And(
			godb.Q("nick = ?", lu.Nick),
			godb.Q("token = ? ", token),
		),
	)
	err := query.Do()
	if err == sql.ErrNoRows {
		return "", lg.New(op, "No user with nick: "+lu.Nick)
	} else if err != nil {
		return "", lg.NewError(op, "Internal error", err)
	}

	return user.Token, nil
}

func (str *Storage) ReadAuthUsers() []auth.AuthUser {
	authUsers := make([]auth.AuthUser, 0)
	users := make([]User, 0)
	_ = str.db.Select(&users).Do()
	for _, u := range users {
		authUsers = append(authUsers, u.ConvertUserToAuth())
	}

	return authUsers
}
