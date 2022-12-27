package memory

import (
	"fmt"

	"github.com/Crandel/go_chat/internal/auth"
	lg "github.com/Crandel/go_chat/internal/logging"
)

func (str *Storage) SigninUser(su auth.SigninUser) (string, error) {
	const op lg.Op = "memory.Signin"
	u := ConvertUserFromSigning(su)
	if str.Users == nil {
		str.Users = make(map[UserId]User)
	}
	str.Lock()
	_, exists := str.Users[u.Nick]
	if exists {
		return "", lg.New(
			op, lg.Info, fmt.Sprintf("User with nick: '%s' exists", u.Nick))
	}
	str.Users[u.Nick] = u
	str.Unlock()
	return u.Token, nil
}

func (str *Storage) LoginUser(lu auth.LoginUser) (string, error) {
	const op lg.Op = "memory.LoginUser"
	str.RLock()
	u, exists := str.Users[UserId(lu.Nick)]
	str.RUnlock()
	if !exists {
		return "", lg.New(
			op, lg.Info, fmt.Sprintf("No user with nick '%s'", lu.Nick))
	}
	token := auth.MakeToken(lu.Nick, lu.Password)
	if u.Token != token {
		return "", lg.New(
			op, lg.Info, fmt.Sprintf("User with nick '%s' has wrong password", lu.Nick))
	}
	return u.Token, nil
}

func (str *Storage) ReadAuthUsers() []auth.AuthUser {
	authUsers := make([]auth.AuthUser, len(str.Users))
	str.RLock()
	for _, u := range str.Users {
		authUsers = append(authUsers, auth.AuthUser{
			Nick:  string(u.Nick),
			Token: u.Token,
		})
	}
	str.RUnlock()
	return authUsers
}
